package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (a *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	records, err := a.DB.GetChirps()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, r := range records {
		chirps = append(chirps, Chirp{
			ID:       r.ID,
			AuthorID: r.AuthorID,
			Body:     r.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondJSON(w, http.StatusOK, chirps)
}

func (a *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	id, err := strconv.Atoi(chirpID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	chirp, err := a.DB.GetChirp(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "couldn't get chirp")
		return
	}

	respondJSON(w, http.StatusOK, Chirp{
		ID:       chirp.ID,
		AuthorID: chirp.AuthorID,
		Body:     chirp.Body,
	})
}
