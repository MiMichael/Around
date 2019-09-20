package main

import (
	"fmt"
	"log"
	"net/http"

	consts "./constant"
	es "./elasticsearch"
	"./rpc"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var (
	mediaTypes = map[string]string{
		".jpeg": "image",
		".jpg":  "image",
		".gif":  "image",
		".png":  "image",
		".mov":  "video",
		".mp4":  "video",
		".avi":  "video",
		".flv":  "video",
		".wmv":  "video",
	}
)

func main() {
	fmt.Println("started-service")
	es.CreateIndexIfNotExist()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(consts.TOKEN_KEY), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r := mux.NewRouter()
	r.Handle("/post", jwtMiddleware.Handler(http.HandlerFunc(rpc.HandlePost))).Methods("POST", "OPTIONS")
	r.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(rpc.HandleSearch))).Methods("GET", "OPTIONS")
	r.Handle("/cluster", jwtMiddleware.Handler(http.HandlerFunc(handlerCluster))).Methods("GET", "OPTIONS")
	r.Handle("/signup", http.HandlerFunc(rpc.HandleRegister)).Methods("POST", "OPTIONS")
	r.Handle("/login", http.HandlerFunc(rpc.HandleLogin)).Methods("POST", "OPTIONS")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
