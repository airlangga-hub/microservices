package main

type Account struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

type Product struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price,omitempty"`
	Quantity    int32  `json:"quantity,omitempty"`
}
