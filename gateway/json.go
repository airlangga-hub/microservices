package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s\n", err)
		http.Error(w, "some error happened", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
