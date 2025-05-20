package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("test_secret_key") // future: load from .env

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string, expirationMinutes int) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	return tokenString, err
}
