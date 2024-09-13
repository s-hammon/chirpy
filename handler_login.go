package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const maxExpire = time.Second * 60 * 60 * 24

func (a *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := a.DB.GetUsers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't login")
		return
	}

	decoder := json.NewDecoder(r.Body)
	body := UserRequest{}
	if err := decoder.Decode(&body); err != nil {
		respondError(w, http.StatusInternalServerError, "couldn't login")
		return
	}

	expire := time.Duration(body.ExpiresInSeconds) * time.Second
	if expire == 0 {
		expire = maxExpire
	}
	fmt.Printf("Token expires after: %v\n", expire)

	for _, user := range dbUsers {
		if body.Email == user.Email {
			if err := bcrypt.CompareHashAndPassword(user.Password, []byte(body.Password)); err == nil {
				claims := &jwt.RegisteredClaims{
					Issuer:    "chirpy",
					IssuedAt:  jwt.NewNumericDate(time.Now().Truncate(time.Millisecond).UTC()),
					ExpiresAt: jwt.NewNumericDate(time.Now().Truncate(time.Millisecond).UTC().Add(expire)),
					Subject:   strconv.Itoa(user.ID),
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				ss, err := token.SignedString([]byte(a.jwtSecret))
				if err != nil {
					respondError(w, http.StatusInternalServerError, "error signing JWT")
				}
				respondJSON(w, http.StatusOK, UserResponse{
					ID:    user.ID,
					Email: user.Email,
					Token: ss,
				})
				return
			}
		}
	}

	respondError(w, http.StatusUnauthorized, "incorrect email or password")
}
