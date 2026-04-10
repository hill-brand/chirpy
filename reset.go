package main

import (
	"context"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "No permission", nil)
		return
	}
	if err := cfg.queries.ResetUsers(context.Background()); err != nil {
		respondWithError(w, 500, "Error resetting server", err)
		return
	}
	cfg.fileserverHits.Store(0)
	type response struct {
		Message string `json:"message"`
	}
	log.Printf("Server reset")
	resp := response{Message: "Server reset successfully"}
	respondWithJSON(w, 200, resp)
}
