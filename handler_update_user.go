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
		fmt.Printf("%v %v", claims.ID, claims.StandardClaims.Issuer)
	} else {
		fmt.Println(err)
	}

	if token.Valid {
		fmt.Println("token valid")
	} else {
		fmt.Println("token not valid")
	}

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
