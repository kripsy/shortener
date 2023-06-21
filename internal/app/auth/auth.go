package auth

import (
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type User struct {
	ID uint64
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

const SECRET_KEY = "supersecretkey"
const TOKEN_EXP = time.Hour * 3

func BuildJWTString() (string, error) {

	userID := rand.Uint64()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", errors.Wrap(err, "failed in BuildJWTString: %w")
	}
	return tokenString, nil
}

func GetUserID(tokenString string) (uint64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
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
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return time.Time{}, err
	}

	if !token.Valid {
		return time.Time{}, err
	}

	return claims.ExpiresAt.Time, nil
}
