package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (a *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := a.DB.GetUsers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't login")
	}
	decoder := json.NewDecoder(r.Body)
	body := UserRequest{}
	if err := decoder.Decode(&body); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't login")
		return
	}

	for _, user := range dbUsers {
		if body.Email == user.Email {
			if err := bcrypt.CompareHashAndPassword(user.Password, []byte(body.Password)); err == nil {
				respondJSON(w, http.StatusOK, UserResponse{
					ID:    user.ID,
					Email: user.Email,
				})
				return
			}
		}
	}

	respondError(w, http.StatusUnauthorized, "incorrect email or password")
}
