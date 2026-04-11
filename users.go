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

// endpoint create a new user
// requires unique email address
func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// decode request
	type newUser struct {
		Email string `json:"email"`
	}
	var user newUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		log.Printf("Error decoding request body: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding request body", err)
		return
	}

	// check for valid request body
	if user.Email == "" {
		log.Printf("Create user request received without valid body")
		respondWithError(w, http.StatusBadRequest, "Include \"email\" in request body", nil)
		return
	}

	// create database entry for user
	result, err := cfg.queries.CreateUser(context.Background(), user.Email)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			respondWithError(w, http.StatusBadRequest, "Email taken, use a different email", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	// return success message
	log.Printf("User created with ID: %v", result.ID)
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
	respondWithJSON(w, http.StatusCreated, response)
}
