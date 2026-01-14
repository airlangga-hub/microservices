package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	accpb "github.com/airlangga-hub/microservices/gateway/account_pb"
)

func (cfg *Config) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := cfg.AccountClient.SignUp(ctx, &accpb.SignUpRequest{Email: request.Email, Password: request.Password})
	if err != nil {
		http.Error(w, "failed to sign up", http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: res.Jwt,
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		response,
	)
}

func (cfg *Config) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := cfg.AccountClient.Login(ctx, &accpb.LoginRequest{Email: request.Email, Password: request.Password})
	if err != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: res.Jwt,
	}

	respondWithJSON(
		w,
		http.StatusOK,
		response,
	)
}
