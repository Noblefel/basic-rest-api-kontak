package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var secret_key = []byte("key")

func GenerateJWT(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
	})

	tokenString, err := token.SignedString(secret_key)
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

func VerifyJWT(s string) (int, error) {
	var claims jwt.MapClaims

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", errors.New("unauthorized")
		}

		return []byte(secret_key), nil
	}

	token, err := jwt.ParseWithClaims(s, &claims, keyFunc)
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("unauthorized")
	}

	userId := claims["user_id"].(float64)
	return int(userId), nil
}
