package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"text/template"
	"time"

	//nolint:gosec
	_ "net/http/pprof"

	//nolint:depguard
	"github.com/kripsy/shortener/internal/app/application"
)

var (
	//nolint:gochecknoglobals
	buildVersion string
	//nolint:gochecknoglobals
	buildDate string
	//nolint:gochecknoglobals
	buildCommit string
)

type BuildData struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

const Template = `	Build version: {{if .BuildVersion}} {{.BuildVersion}} {{else}} N/A {{end}}
	Build version: {{if .BuildDate}} {{.BuildDate}} {{else}} N/A {{end}}
	Build version: {{if .BuildCommit}} {{.BuildCommit}} {{else}} N/A {{end}}
`

func main() {
	ctx := context.Background()

	application, err := application.NewApp(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}

	d := &BuildData{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}

	t := template.Must(template.New("buildTags").Parse(Template))

	err = t.Execute(os.Stdout, *d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}
	defer func() { // flushes buffer, if any
		if err = application.GetAppLogger().Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)

			return
		}
	}()

	defer application.GetAppRepo().Close() // close repo

	fmt.Printf("SERVER_ADDRESS: %s\n", application.GetAppConfig().URLServer)
	fmt.Printf("BASE_URL: %s\n", application.GetAppConfig().URLPrefixRepo)

	srv := &http.Server{
		Addr:         application.GetAppConfig().URLServer,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		//nolint:gomnd
		IdleTimeout: 30 * time.Second,
		//nolint:gomnd
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           application.GetAppServer().Router,
	}
	if application.GetAppConfig().EnableHTTPS != "" {
		fmt.Println("CREATE CERT")
		err = createCertificate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)

			return
		}
		fmt.Println("Success CREATE CERT")
		err = srv.ListenAndServeTLS("./pki/server.crt", "./pki/server.key")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)

			return
		}

		return
	}
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}
}

func createCertificate() error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"EngeniyOrg"},
			Country:      []string{"RU"},
		},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)},
		NotBefore:   time.Now(),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
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
		return err
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return fmt.Errorf("error encode private key %w", err)
	}

	err = saveCert("./pki/server.crt", &certPEM)
	if err != nil {
		return err
	}
	err = saveCert("./pki/server.key", &privateKeyPEM)
	if err != nil {
		return err
	}
	return nil
}

func saveCert(path string, payload *bytes.Buffer) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("error open file %w", err)
	}
	writer := bufio.NewWriter(f)
	_, err = writer.ReadFrom(payload)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
