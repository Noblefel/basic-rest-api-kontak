package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var secret_key = []byte("key")

func GenerateJWT(userId int) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret_key)
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

func VerifyJWT(s string) (float64, error) {
	var userId float64

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(s, claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", errors.New("Unauthorized")
		}

		return []byte(secret_key), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("Unauthorized")
	}

	userId = claims["user_id"].(float64)

	return userId, nil
}
