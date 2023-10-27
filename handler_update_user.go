package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type jwtClaim struct {
	Sub    int `json:"sub"`
	Expiry int `json:"expires"`
	jwt.StandardClaims
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	tokenStr := r.Header.Get("Authorization")
	secretKey := []byte(jwtSecret)
	fmt.Println(tokenStr)

	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(*jwtClaim); ok && token.Valid {
		user, err := cfg.DB.UpdateUser(params.Email, params.Password, claims.Sub)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
			return
		}
		fmt.Println(claims.Sub)
		respondWithJSON(w, http.StatusOK, User{
			ID:    claims.Sub,
			Email: user.Email,
		})
	} else {
		fmt.Println(err)
	}

}

func handlerUpdateUser(jwt string) bool {
	return false
}
