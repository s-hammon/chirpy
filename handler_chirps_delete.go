package main

import (
	"net/http"
	"strconv"

	"github.com/s-hammon/chirpy/internal/auth"
)

func (a *apiConfig) handleDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "couldn't parse JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, a.jwtSecret)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "couldn't authorize JWT")
		return
	}
	authorID, err := strconv.Atoi(subject)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't parse user ID")
		return
	}

	chirpID := r.PathValue("chirpID")
	id, err := strconv.Atoi(chirpID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	if err := a.DB.DeleteChirpByID(id, authorID); err != nil {
		if err.Error() == "permission denied" {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusNotFound, "couldn't delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
