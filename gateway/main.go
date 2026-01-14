package main

import "os"

var jwtSecret []byte

func main() {
	secret := os.Getenv("JWT_SECRET")
	jwtSecret = []byte(secret)
	
	// public endpoints
	// buyer endpoints
	// seller endpoints
}
