package main

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func MakeJWT(id int32, email string, secret []byte) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": id,
			"user_email": email,
			"exp": time.Now().Add(time.Hour * 24).Format(time.DateOnly),
		},
	)
	
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	
	return tokenStr, nil
}