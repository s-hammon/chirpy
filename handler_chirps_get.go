package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (a *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := a.DB.GetChirps()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []ChirpClean{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, ChirpClean{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondJSON(w, http.StatusOK, chirps)
}

func (a *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := a.DB.GetChirps()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't retrieve chirps")
		return
	}

	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondError(w, http.StatusBadRequest, "must provide id")
		return
	}
	id, err := strconv.Atoi(chirpID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "id must be an integer")
	}

	chirp := ChirpClean{}
	for _, ch := range dbChirps {
		if ch.ID == id {
			chirp.ID = ch.ID
			chirp.Body = ch.Body
		}

	}

	if chirp.ID != 0 {
		respondJSON(w, http.StatusOK, chirp)
	} else {
		respondError(w, http.StatusNotFound, "couldn't find chirp")
	}
}
