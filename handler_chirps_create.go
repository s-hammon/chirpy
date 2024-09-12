package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ChirpBody struct {
	Body string `json:"body"`
}

type ChirpError struct {
	Error string `json:"error"`
}

type ChirpAck struct {
	Valid bool `json:"valid"`
}

type ChirpClean struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (a *apiConfig) handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := ChirpBody{}
	if err := decoder.Decode(&body); err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't read parameters")
		return
	}

	const maxLength = 140
	if len(body.Body) > maxLength {
		respondError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleanedChirp := censorProfanity(body.Body)

	chirp, err := a.DB.CreateChirp(cleanedChirp)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondJSON(w, http.StatusCreated, ChirpClean{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}

func censorProfanity(msg string) string {
	naughtyWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(msg, " ")
	for i, word := range words {
		lowered := strings.ToLower(word)
		if _, ok := naughtyWords[lowered]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
