package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/s-hammon/chirpy/internal/auth"
)

const maxExpire = time.Second * 60 * 60 * 24

func (a *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't login")
		return
	}

	user, err := a.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	if err = auth.CheckHash(user.Password, params.Password); err != nil {
		respondError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	expire := time.Duration(params.ExpiresInSeconds) * time.Second
	if expire == 0 {
		expire = maxExpire
	}

	token, err := auth.MakeJWT(user.ID, a.jwtSecret, expire)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create refresh token")
		return
	}

	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 60)
	if err = a.DB.CreateRefreshToken(user.ID, refreshToken, expiresAt); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't write refresh token")
	}

	respondJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
