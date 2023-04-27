package main

import (
	"log"
	"net/http"
)


func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./fileserver/"))
	
	mux.Handle("/fileserver/", http.StripPrefix("/fileserver", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/password/create", passwordCreate)
	mux.HandleFunc("/password/view", passwordView)

	log.Print("Listening on port :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
