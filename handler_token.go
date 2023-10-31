package main

import (
	"RossDooney/go-webserver-test/internal/auth"
	"fmt"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, issuer, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	if issuer != "chirpy-refresh" {
		fmt.Println("not access token")
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	fmt.Println("refresh token submitted")
	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	fmt.Printf("user id: %v \n", userID)

	newToken, err := auth.MakeJWT(userID, cfg.jwtSecret, "chirpy-access", 3600)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't refresh JWT")
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	err = cfg.DB.CreateRevokeToken(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token")
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}
