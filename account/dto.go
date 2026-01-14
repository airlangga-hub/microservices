package main

type Account struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password,omitempty"`
	Type           string `json:"type"`
}
