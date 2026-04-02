package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	bannedWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanMessage(params.Body, bannedWords),
	})
}

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
