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

type validLogin struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type login struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Expires  int    `json:"expires_in_seconds"`
}

func (cfg *apiConfig) handlerVerifyLogin(w http.ResponseWriter, r *http.Request) {
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
	jwt, err := CreateJwt(user.ID)
	if err != nil {
		log.Fatal(err)
		return
	}
	respondWithJSON(w, http.StatusOK, validLogin{
		Email: user.Email,
		ID:    user.ID,
		Token: jwt,
	})
}

func CreateJwt(id int) (string, error) {
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["sub"] = id
	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println(tokenStr)
	return tokenStr, nil
}
