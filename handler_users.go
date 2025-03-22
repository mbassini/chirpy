package main

import (
	"encoding/json"
	"net/http"

	"github.com/lib/pq"
	"github.com/mbassini/chirpy/internal/database"
)

func handlerCreateUser(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type createUserRequest struct {
			Email string `json:"email"`
		}

		var req createUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}

		user, err := db.CreateUser(r.Context(), req.Email)
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

		respondWithJSON(w, http.StatusCreated, user)
	}
}
