package main

import (
	"net/http"
	"os"
)

var jwtSecret []byte

func main() {
	secret := os.Getenv("JWT_SECRET")
	jwtSecret = []byte(secret)

	// public endpoints
	http.HandleFunc("POST /api/signup", SignUp)
	// buyer endpoints
	// seller endpoints
}
