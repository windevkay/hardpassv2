package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/password/create", passwordCreate)
	mux.HandleFunc("/password/view", passwordView)

	log.Print("Listening on port :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
