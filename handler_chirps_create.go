package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/s-hammon/chirpy/internal/auth"
)

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (a *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	subject, err := auth.ValidateJWT(token, a.jwtSecret)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "couldn't validate JWT")
		return
	}

	authorID, err := strconv.Atoi(subject)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "couldn't extract user ID from JWT")
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't read parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := a.DB.CreateChirp(authorID, cleaned)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}

	respondJSON(w, http.StatusCreated, Chirp{
		ID:       chirp.ID,
		AuthorID: chirp.AuthorID,
		Body:     chirp.Body,
	})
}

func validateChirp(body string) (string, error) {
	const maxLength = 140
	if len(body) > maxLength {
		return "", errors.New("chirp is too long")
	}

	naughtyWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := censorProfanity(body, naughtyWords)
	return cleaned, nil
}

func censorProfanity(msg string, naughtyWords map[string]struct{}) string {
	words := strings.Split(msg, " ")
	for i, word := range words {
		lowered := strings.ToLower(word)
		if _, ok := naughtyWords[lowered]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
