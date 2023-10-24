package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerVerifyLogin(w http.ResponseWriter, r *http.Request) {
	type login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := login{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CheckUserLogin(params.Email, params.Password)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't login")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		Email: user.Email,
		ID:    user.ID,
	})
}
