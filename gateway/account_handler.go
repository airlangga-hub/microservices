package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	accpb "github.com/airlangga-hub/microservices/gateway/account_pb"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.AccountClient.SignUp(ctx, &accpb.SignUpRequest{Email: request.Email, Password: request.Password})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		struct {
			Token string `json:"token"`
		}{
			Token: res.Jwt,
		},
	)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.AccountClient.Login(ctx, &accpb.LoginRequest{Email: request.Email, Password: request.Password})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		struct {
			Token string `json:"token"`
		}{
			Token: res.Jwt,
		},
	)
}
