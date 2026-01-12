package main

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func MakeJWT(id int32, email string, secret []byte) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * 24),
			"sub": email,
			"id":  id,
		},
	)

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func VerifyJWT(authHeader string, secret []byte) (Account, error) {
	tokenSlice := strings.Fields(authHeader)

	if len(tokenSlice) < 2 || tokenSlice[0] != "Bearer" {
		return Account{}, errors.New("malformed auth header")
	}

	token, err := jwt.Parse(
		tokenSlice[1],
		func(t *jwt.Token) (any, error) {
			return secret, nil
		},
	)
	if err != nil {
		return Account{}, err
	}
	if !token.Valid {
		return Account{}, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Account{}, errors.New("invalid token claims type")
	}

	id, ok := claims["id"].(int32)
	if !ok {
		return Account{}, errors.New("invalid token")
	}

	email, ok := claims["sub"].(string)
	if !ok {
		return Account{}, errors.New("invalid token")
	}

	return Account{
		ID:    id,
		Email: email,
	}, nil
}
