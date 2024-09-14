package main

import (
	"encoding/json"
	"net/http"

	"github.com/s-hammon/chirpy/internal/auth"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (a *apiConfig) handleNewUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	pwd, err := auth.HashPasword(params.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	user, err := a.DB.CreateUser(params.Email, pwd)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	respondJSON(w, http.StatusCreated, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
