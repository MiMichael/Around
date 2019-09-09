package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	ex "../external"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one register request")
	w.Header().Set("Content-Type", "text/plain")
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

	if user.Username == "" || user.Password == "" || !regexp.MustCompile(`^[a-z0-9_]+$`).MatchString(user.Username) {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		fmt.Printf("Invalid username or password.\n")
		return
	}

	if err := ex.AddUser(user); err != nil {
		if err.Error() == "User already exists" {
			http.Error(w, "User already exists", http.StatusBadRequest)
			fmt.Printf("User already exists: %v.\n", err)
		} else {
			http.Error(w, "Failed to save to ElasticSearch", http.StatusInternalServerError)
			fmt.Printf("Failed to save to ElasticSearch: %v.\n", err)
		}
		return
	}
	w.Write([]byte("User added successfully"))

}
