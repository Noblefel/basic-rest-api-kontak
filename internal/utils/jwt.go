package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var secret_key = []byte("key")

func GenerateJWT(userId, level int) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["level"] = level

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret_key)
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

func VerifyJWT(s string) (float64, float64, error) {
	var userId, level float64

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(s, claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", errors.New("Unauthorized")
		}

		return []byte(secret_key), nil
	})

	if err != nil {
		return 0, 0, err
	}

	if !token.Valid {
		return 0, 0, errors.New("Unauthorized")
	}

	userId = claims["user_id"].(float64)

	if claims["level"] == nil {
		level = 0
	} else {
		level = claims["level"].(float64)
	}

	return userId, level, nil
}
