package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var TokenDuration = time.Hour * 24 * 7

type Claims struct {
	jwt.StandardClaims
	ID   int    `json:"id"`
	Role string `json:"role"`
	Exp  int64  `json:"exp"`
}

func GenerateToken(id int, tokenDuration time.Duration, secretKey string) (string, error) {
	claims := Claims{
		ID:  id,
		Exp: time.Now().Add(tokenDuration * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secretKey string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims := token.Claims.(*Claims)
	return claims, nil
}

