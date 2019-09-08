package main

import (
	"fmt"
	"log"
	"net/http"

	es "./elasticsearch"
	"./rpc"
)

func main() {
	fmt.Println("started-service")
	es.CreateIndexIfNotExist()

	http.HandleFunc("/post", rpc.HandlePost)
	http.HandleFunc("/search", rpc.HandleSearch)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
