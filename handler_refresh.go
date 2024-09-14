package main

import (
	"net/http"
	"time"

	"github.com/s-hammon/chirpy/internal/auth"
)

func (a *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshToken, err := a.DB.GetRefreshTokenByValue(authToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID := refreshToken.UserID
	expiresIn := time.Second * 60 * 60
	token, err := auth.MakeJWT(userID, a.jwtSecret, expiresIn)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create JWT")
	}

	respondJSON(w, http.StatusOK, response{
		Token: token,
	})
}
