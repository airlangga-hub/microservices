package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func MakeJWT(userID int32, email string, secret []byte) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp":       time.Now().UTC().Add(time.Hour * 24).Unix(),
			"sub":       email,
			"user_id":   userID,
		},
	)

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
