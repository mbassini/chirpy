package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
	}
	type validResponse struct {
		Valid bool `json:"valid"`
	}

	var req chirpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	const maxChirpLength = 140
	if len(req.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, validResponse{Valid: true})
}
