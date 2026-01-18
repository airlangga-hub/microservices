package main

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func MakeJWTSeller(userID int32, email string, secret []byte) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp":       time.Now().UTC().Add(time.Hour * 24).Unix(),
			"sub":       email,
			"user_type": "seller",
			"user_id":   userID,
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

	userType, ok := claims["user_type"].(string)
	if !ok {
		log.Println("ERROR: gateway VerifyJWT: parsing claims user type")
		return Account{}, errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		log.Println("ERROR: gateway VerifyJWT: parsing claims user id")
		return Account{}, errors.New("invalid token")
	}

	email, ok := claims["sub"].(string)
	if !ok {
		log.Println("ERROR: gateway VerifyJWT: parsing claims sub")
		return Account{}, errors.New("invalid token")
	}

	return Account{
		ID:    int32(userID),
		Email: email,
		Type:  userType,
	}, nil
}
