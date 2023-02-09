package main

import (
	"log"
	"net/http"
	handler "smolage/images/api"
)

func main() {
	http.HandleFunc("/api/smol/", handler.Handler)
	log.Printf("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
