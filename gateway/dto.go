package main

type Account struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}
