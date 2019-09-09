package rpc

import (
	"fmt"
	"net/http"
	"strconv"

	consts "../constant"
	es "../elasticsearch"
	ex "../external"
	gcs "../gcs"
	jwt "github.com/dgrijalva/jwt-go"
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

	user := r.Context().Value("user")
	claims := user.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]

	lat, _ := strconv.ParseFloat(r.FormValue("lat"), 64)
	lon, _ := strconv.ParseFloat(r.FormValue("lon"), 64)

	p := &ex.Post{
		User:    username.(string),
		Message: r.FormValue("message"),
		Location: ex.Location{
			Lat: lat,
			Lon: lon,
		},
	}

	id := uuid.New().String()

	image, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image is not available", http.StatusBadRequest)
		fmt.Printf("Image is not available %v\n", err)
		return
	}
	attrs, err := gcs.SaveToGCS(image, consts.BUCKET_NAME, id)
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
