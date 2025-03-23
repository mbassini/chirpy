package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mbassini/chirpy/internal/database"
)

const errorForeignKeyViolation = "23503"

func handlerCreateChirp(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type createChirpRequest struct {
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}

		var req createChirpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusInternalServerError, "something went wrong", err)
			return
		}

		chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: req.Body, UserID: req.UserID})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == errorForeignKeyViolation {
					respondWithError(w, http.StatusNotFound, "user does not exist", err)
					return
				}
			}
			respondWithError(w, 400, "error creating chirp", err)
			return
		}

		myChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}

		respondWithJSON(w, http.StatusCreated, myChirp)
	}
}
