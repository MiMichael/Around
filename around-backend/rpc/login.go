package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	consts "../constant"
	ex "../external"
	jwt "github.com/dgrijalva/jwt-go"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one login request")
	w.Header().Set("Content-type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var user ex.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode user data from client: %v.\n", err)
		return
	}
	if err := ex.CheckUser(user.Username, user.Password); err != nil {
		if err.Error() == "Wrong username or password" {
			http.Error(w, "Wrong username or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Failed to read from ElasticSearch", http.StatusInternalServerError)
		}
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(consts.TOKEN_KEY))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		fmt.Printf("Failed to generate token %v.\n", err)
		return
	}
	w.Write([]byte(tokenString))

}
