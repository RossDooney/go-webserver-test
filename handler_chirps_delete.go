package main

import (
	"RossDooney/go-webserver-test/internal/auth"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
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
	if issuer != "chirpy-access" {
		fmt.Println("not access token")
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	userIDInt, _ := strconv.Atoi(subject)

	if userIDInt != chirpID {
		respondWithError(w, http.StatusForbidden, "Not Authorized")
		return
	}

	cfg.DB.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't delete chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{})
}
