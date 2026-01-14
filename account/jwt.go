package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func MakeJWT(userType, email string, secret []byte) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp":       time.Now().Add(time.Hour * 24),
			"sub":       email,
			"user_type": userType,
		},
	)

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

