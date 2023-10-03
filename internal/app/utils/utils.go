// Package utils provides the helpful functionality for shortener.
package utils

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/fs"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	//nolint:depguard

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

const (
	ServerCertPath = "./cert/server.crt"
	PrivateKeyPath = "./cert/server.key"
)

// CreateShortURL returns random short URL, consist 5 bytes (Why not?).
// If we have same error, it returns empty string and error.
func CreateShortURL() (string, error) {
	reqLen := 5
	buf := make([]byte, reqLen)
	_, err := rand.Read(buf)

	if err != nil {
		//nolint:goerr113,nolintlint
		return "", fmt.Errorf("error while generating random string: %w", err)
	}

	return fmt.Sprintf("%x", buf), nil
}

// CreateShortURLWithoutFmt returns random short URL, consist 5 bytes.
// It's optimized version of CreateShortURL.
// If we have same error, it returns empty string and error.
func CreateShortURLWithoutFmt() (string, error) {
	reqLen := 5
	buf := make([]byte, reqLen)
	_, err := rand.Read(buf)

	if err != nil {
		//nolint:goerr113,nolintlint
		return "", fmt.Errorf("error while generating random string: %w", err)
	}

	return hex.EncodeToString(buf), nil
}

// ReturnURL returns an union shortURL and address of our server.
func ReturnURL(endpoint, globalURL string) string {
	return globalURL + "/" + endpoint
}

// StingContains returns is searchString contain in arrayString.
func StingContains(arrayString []string, searchString string) bool {
	for _, v := range arrayString {
		if v == searchString {
			return true
		}
	}

	return false
}

// GetTokenFromBearer returns token from header.
// Header should start from "Baerer ", otherwise return empty string and error.
func GetTokenFromBearer(bearerString string) (string, error) {
	reqLen := 2
	splitString := strings.Split(bearerString, "Bearer ")
	fmt.Printf("len splitString %d\n", len(splitString))
	if len(splitString) < reqLen {
		fmt.Printf("bearer string not valid")
		//nolint:goerr113
		return "", fmt.Errorf("bearer string not valid")
	}
	tokenString := splitString[1]
	if tokenString == "" {
		fmt.Printf("tokenString is empty")

		//nolint:goerr113
		return "", fmt.Errorf("tokenString is empty")
	}

	return tokenString, nil
}

// GetToken returns token from header or cookie.
// Header should start from "Baerer ", otherwise return empty string and error.
func GetToken(r *http.Request) (string, error) {
	var token string

	tokenString := r.Header.Get("Authorization")
	if tokenString != "" {
		fmt.Printf("get token from header: %s\n", tokenString)
		token, _ = GetTokenFromBearer(tokenString)
		fmt.Printf("token %s\n", token)
	}
	if token != "" {
		return token, nil
	}

	// if we continue - it means that in header isn't token. Try find it in cookie
	cookieToken, err := r.Cookie("token")
	if err != nil {
		return "", errors.Wrap(err, "cannot get token from cookie")
	}
	token = cookieToken.Value
	if token != "" {
		return token, nil
	}
	//nolint:goerr113
	return "", fmt.Errorf("token not found")
}

func CreateCertificate() error {
	maxInt := 1658
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(int64(maxInt)),
		Subject: pkix.Name{
			Organization: []string{"EngeniyOrg"},
			Country:      []string{"RU"},
		},
		//nolint:gomnd
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)},
		NotBefore:   time.Now(),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}
	//nolint:gomnd
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("error generate key %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("error create certificate %w", err)
	}

	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return fmt.Errorf("error encode cert %w", err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return fmt.Errorf("error encode private key %w", err)
	}

	err = saveCert(ServerCertPath, &certPEM)
	if err != nil {
		return err
	}
	err = saveCert(PrivateKeyPath, &privateKeyPEM)
	if err != nil {
		return err
	}

	return nil
}

func saveCert(path string, payload *bytes.Buffer) error {
	permissionValue := 0755
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fs.FileMode(permissionValue))
	if err != nil {
		return fmt.Errorf("error open file %w", err)
	}
	writer := bufio.NewWriter(f)
	_, err = writer.ReadFrom(payload)
	if err != nil {
		return fmt.Errorf("error write to file %w", err)
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("error close file %w", err)
	}

	return nil
}

func GetTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {

		return "", fmt.Errorf("not metadata in context")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", fmt.Errorf("token not found")
	}

	token := strings.TrimPrefix(values[0], "Bearer ")
	return token, nil
}
