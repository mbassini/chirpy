package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const errorUniqueViolation = "23505"

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Email string `json:"email"`
	}

	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameteres", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == errorUniqueViolation {
				respondWithError(w, http.StatusConflict, "Email already in use", err)
				return
			}
		}
		respondWithError(w, 400, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

func (cfg *apiConfig) handlerUsersReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "reset endpoint is only available in development mode", nil)
		return
	}

	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not reset users", err)
		return
	}

	respondWithJSON(w, http.StatusOK, "Users table has been reset")
}
