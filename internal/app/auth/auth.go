// Package auth provides functionality for wotking with jwt token.
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const secretKey = "supersecretkey"
const tokenExp = time.Hour * 3

func BuildJWTString(userID int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.Wrap(err, "failed in BuildJWTString: %w")
	}
	return tokenString, nil
}

func GetUserID(tokenString string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, err
	}

	return claims.UserID, nil
}

func GetExpires(tokenString string) (time.Time, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return time.Time{}, err
	}

	if !token.Valid {
		return time.Time{}, err
	}

	return claims.ExpiresAt.Time, nil
}

func IsTokenValid(tokenString string) (bool, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, err
	}

	return true, nil
}
