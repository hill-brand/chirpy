package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type newUser struct {
		Email string `json:"email"`
	}
	var user newUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		log.Printf("Error decoding request body: %v", err)
		respondWithError(w, 500, "Error decoding request body", err)
		return
	}
	if user.Email == "" {
		log.Printf("Create user request received without valid body")
		respondWithError(w, 400, "Include \"email\" in request body", nil)
		return
	}
	result, err := cfg.queries.CreateUser(context.Background(), user.Email)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			respondWithError(w, 400, "Email taken, use a different email", err)
			return
		}
		respondWithError(w, 500, "Error creating user", err)
		return
	}
	type Response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	response := Response{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Email:     result.Email,
	}
	respondWithJSON(w, 201, response)
}