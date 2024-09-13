package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (a *apiConfig) handleNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := User{}
	if err := decoder.Decode(&body); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	user, err := a.DB.CreateUser(body.Email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	respondJSON(w, http.StatusCreated, User{
		ID:    user.ID,
		Email: user.Email,
	})
}
