package main

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

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

	userType, ok := claims["user_type"].(string)
	if !ok {
		return Account{}, errors.New("invalid token")
	}

	email, ok := claims["sub"].(string)
	if !ok {
		return Account{}, errors.New("invalid token")
	}

	return Account{
		Email: email,
		Type:  userType,
	}, nil
}
