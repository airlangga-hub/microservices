package main

import (
	"context"
	"net/http"
)

type contextKey string

const accountKey contextKey = "account"

func AuthorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		account, err := VerifyJWT(authHeader, jwtSecret)
		if err != nil {
			http.Error(w, "unauthorized user", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), accountKey, account)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
