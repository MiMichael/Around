package rpc

import (
	"fmt"
	"net/http"
	"strconv"

	consts "../constant"

	es "../elasticsearch"
	gcs "../gcs"
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
	lat, _ := strconv.ParseFloat(r.FormValue("lat"), 64)
	lon, _ := strconv.ParseFloat(r.FormValue("lon"), 64)

	p := &post.Post{
		User:    r.FormValue("user"),
		Message: r.FormValue("message"),
		Location: post.Location{
			Lat: lat,
			Lon: lon,
		},
	}

	id := uuid.New().String()

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image is not available", http.StatusBadRequest)
		fmt.Printf("Image is not available %v\n", err)
		return
	}
	attrs, err := gcs.SaveToGCS(file, consts.BUCKET_NAME, id)
	if err != nil {
		http.Error(w, "Failed to save image to GCS", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to GCS %v.\n", err)
		return

	}
	p.Url = attrs.MediaLink

	err = es.SaveToES(p, id)
	if err != nil {
		http.Error(w, "Failed to save post to Elastic Search", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to Elastic Search %v\n", err)
		return

	}
	fmt.Printf("Saved one post to Elastic Search: %s", p.Message)
}
