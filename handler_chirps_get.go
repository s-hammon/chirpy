package main

import (
	"net/http"
	"sort"
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
