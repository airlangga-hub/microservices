package main

type Account struct {
	ID             int32  `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}
