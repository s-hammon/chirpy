package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (a *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := -1
	reqAuthorId := r.URL.Query().Get("author_id")
	if reqAuthorId != "" {
		intID, err := strconv.Atoi(reqAuthorId)
		if err != nil {
			respondError(w, http.StatusBadRequest, "author_id must be an integer")
			return
		}

		authorID = intID
	}

	sortQuery := r.URL.Query().Get("sort")
	if sortQuery == "" {
		sortQuery = "asc"
	}

	records, err := a.DB.GetChirps(authorID)
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
		if sortQuery == "desc" {
			return chirps[i].ID > chirps[j].ID
		}
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
