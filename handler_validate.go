package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
	}
	type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
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

	cleanChirp(&req.Body)

	respondWithJSON(w, http.StatusOK, validResponse{CleanedBody: req.Body})
}

func cleanChirp(s *string) {
	bannedWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(*s, " ")
	for i, word := range words {
		_, isBanned := bannedWords[strings.ToLower(word)]
		if isBanned {
			words[i] = "****"
		}
	}

	*s = strings.Join(words, " ")
}
