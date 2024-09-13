package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (a *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := a.DB.GetUsers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't login")
		return
	}

	auth := r.Header.Get("Authorization")

	if auth == "" {
		respondError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	tokenString, ok := strings.CutPrefix(auth, "Bearer ")
	if !ok {
		respondError(w, http.StatusUnauthorized, "invalid token header")
	}

	claims := &jwt.RegisteredClaims{}
	if _, err = jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.jwtSecret), nil
		}); err != nil {
		respondError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	if claims.Issuer != "chirpy" {
		respondError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	if claims.ExpiresAt.UTC().Before(time.Now().Truncate(time.Millisecond).UTC()) {
		respondError(w, http.StatusUnauthorized, "token has expired")
	}

	userId, err := claims.GetSubject()
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	id, _ := strconv.Atoi(userId)

	decoder := json.NewDecoder(r.Body)
	body := UserRequest{}
	if err := decoder.Decode(&body); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't parse request")
		return
	}

	for _, user := range dbUsers {
		if id == user.ID {
			pwd, err := bcrypt.GenerateFromPassword([]byte(body.Password), 1)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "couldn't update user")
				return
			}
			updated, err := a.DB.UpdateUser(id, body.Email, pwd)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "couldn't update user")
				return
			}

			respondJSON(w, http.StatusOK, &UserResponse{
				ID:    id,
				Email: updated.Email,
			})
		}
	}

}
