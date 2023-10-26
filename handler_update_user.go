package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type jwtClaim struct {
	ID     int `json:"id"`
	Expiry int `json:"expires"`
	jwt.StandardClaims
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	jwt := r.Header.Get("Authorization")
	fmt.Println(r.Header.Get(jwt))

	user, err := cfg.DB.UpdateUser(params.Email, params.Password, params.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:    user.ID,
		Email: user.Email,
	})
}

func handlerUpdateUser(jwt string) bool {
	return false
}
