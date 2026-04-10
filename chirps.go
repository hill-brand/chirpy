package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/luckysal/chirpy/internal/database"
)

// post a chirp to the database
// requires json body and user_id
func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	// constants
	const MAX_CHIRP_LENGTH = 140
	BANNED_WORDS := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	// decode request body
	type newChirp struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	var chirp newChirp
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&chirp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	// validate request body
	if chirp.Body == "" || chirp.UserID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Chirps require a body and a user_id", nil)
		return
	}
	if len(chirp.Body) > MAX_CHIRP_LENGTH {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// clean body
	cleanedBody := cleanMessage(chirp.Body, BANNED_WORDS)
	params := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: chirp.UserID,
	}

	// post to database
	result, err := cfg.queries.CreateChirp(context.Background(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	// success message
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	log.Printf("Chirp posted with id: %v", result.ID)
	resp := response{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Body:      result.Body,
		UserID:    result.UserID,
	}
	respondWithJSON(w, http.StatusCreated, resp)
}

// finds banned words in message
// replaces banned words with "****"
// returns new message
func cleanMessage(message string, bannedWords map[string]struct{}) string {
	words := strings.Split(message, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := bannedWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	newMessage := strings.Join(words, " ")
	return newMessage
}
