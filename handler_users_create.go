package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (a *apiConfig) handleNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := UserRequest{}
	if err := decoder.Decode(&body); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(body.Password), 1)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	user, err := a.DB.CreateUser(body.Email, pwd)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	respondJSON(w, http.StatusCreated, UserResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}
