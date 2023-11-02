package main

import (
	"RossDooney/go-webserver-test/internal/auth"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID int `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}
	token, err := auth.GetApikey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find Apikey")
		return
	}
	if token != cfg.apiKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}
	err = cfg.DB.UpdateUserRed(params.Data.UserID, true)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}
	w.WriteHeader(http.StatusOK)
}
