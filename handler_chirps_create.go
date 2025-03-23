package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mbassini/chirpy/internal/database"
)

const errorForeignKeyViolation = "23503"

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type createChirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	var req createChirpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	cleanedChirp, err := validateChirp(req.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanedChirp, UserID: req.UserID})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == errorForeignKeyViolation {
				respondWithError(w, http.StatusNotFound, "User does not exist", err)
				return
			}
		}
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	bannedWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := sanitizeBody(body, bannedWords)
	return cleaned, nil
}

func sanitizeBody(body string, bannedWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		_, isBanned := bannedWords[strings.ToLower(word)]
		if isBanned {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
