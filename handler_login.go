package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func (cfg *apiConfig) handlerVerifyLogin(w http.ResponseWriter, r *http.Request) {
	type login struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Expires  int    `json:"expires_in_seconds"`
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
	jwt, err := CreateJwt()
	if err != nil {
		log.Fatal(err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		Email: user.Email,
		ID:    user.ID,
	})
}

func CreateJwt() (string, error) {
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	tokenStr, err := token.SignedString(jwtSecret)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return tokenStr, nil
}
