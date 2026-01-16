package main

import (
	"context"
	"net/http"
)

type contextKey string

const accountKey contextKey = "account"

func (h *Handler) AuthorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		account, err := VerifyJWT(authHeader, h.Secret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), accountKey, account)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) AuthorizeSellerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		account, err := VerifyJWT(authHeader, h.Secret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		
		if account.Type != "seller" {
			http.Error(w, "you are not a seller", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), accountKey, account)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}