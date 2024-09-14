package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/s-hammon/chirpy/internal/auth"
)

func (a *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}
	token, err := auth.GetToken("Bearer", r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	subject, err := auth.ValidateJWT(token, a.jwtSecret)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "couldn't validate JWT")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't parse request")
		return
	}

	hashedPwd, err := auth.HashPasword(params.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't parse user ID")
		return
	}

	user, err := a.DB.UpdateUser(userID, params.Email, hashedPwd)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't update user")
		return
	}

	respondJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}
