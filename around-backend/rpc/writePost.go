package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	es "../elasticsearch"
	"../post"
	"github.com/google/uuid"
)

func HandlePost(w http.ResponseWriter, r *http.Request) {
	// parse from body of request to get a json object.
	fmt.Println("Received one post request")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	if r.Method == "OPTIONS" {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var p post.Post
	if err := decoder.Decode(&p); err != nil {
		http.Error(w, "Cannot decode post data from client", http.StatusBadRequest)
		fmt.Printf("Cannot decode post data from client %v\n", err)
		return
	}
	id := uuid.New().String()
	err := es.SaveToES(&p, id)
	if err != nil {
		http.Error(w, "Failed to save post to Elastic Search", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to Elastic Search %v\n", err)
		return

	}
	fmt.Printf("Saved one post to Elastic Search: %s", p.Message)
}
