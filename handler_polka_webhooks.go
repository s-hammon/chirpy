package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/s-hammon/chirpy/internal/auth"
)

const userUpgraded = "user.upgraded"

type Data struct {
	UserID int `json:"user_id"`
}

func (a *apiConfig) handlePolkaWebhookUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  Data
	}

	token, err := auth.GetToken("ApiKey", r.Header)
	if err != nil || token != a.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if params.Event != userUpgraded {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := a.DB.UpdateUserUpgrade(params.Data.UserID); err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
