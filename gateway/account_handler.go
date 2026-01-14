package main

import (
	"net/http"
)

func (cfg *Config) SignUp(w http.ResponseWriter, r *http.Request) {
	res, err := accountClient.SignUp()
}