package main

import (
	"encoding/json"
	"net/http"

	"github.com/lib/pq"
)

func handlerCreateUser(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type createUserRequest struct {
			Email string `json:"email"`
		}

		var req createUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}

		user, err := cfg.db.CreateUser(r.Context(), req.Email)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" {
					respondWithError(w, http.StatusConflict, "Email already in use", err)
					return
				}
			}
			respondWithError(w, 400, "Error creating user", err)
			return
		}

		myUser := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		respondWithJSON(w, http.StatusCreated, myUser)
	}
}

func handlerResetUsers(cfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
