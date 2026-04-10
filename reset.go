// testing
package main

import (
	"context"
	"log"
	"net/http"
)

// endpoint reset database and server hits metric
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	// check for dev environment
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "No permission", nil)
		return
	}

	// reset database
	if err := cfg.queries.ResetUsers(context.Background()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error resetting server", err)
		return
	}

	// reset server hits metric
	cfg.fileserverHits.Store(0)

	// return success message
	type response struct {
		Message string `json:"message"`
	}
	log.Printf("Server reset")
	resp := response{Message: "Server reset successfully"}
	respondWithJSON(w, http.StatusOK, resp)
}
