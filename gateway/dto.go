package main

type Account struct {
	ID    int32  `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Type  string `json:"type,omitempty"`
}

type Product struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price,omitempty"`
	Quantity    int32  `json:"quantity,omitempty"`
}
