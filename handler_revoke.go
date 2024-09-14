package main

import (
	"net/http"

	"github.com/s-hammon/chirpy/internal/auth"
)

func (a *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := a.DB.DeleteRefreshTokenByValue(authToken); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
