package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hardpassv2"))
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Uploading file"))
}

func viewFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Viewing file"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/file/upload", uploadFile)
	mux.HandleFunc("/file/view", viewFile)

	log.Print("Listening on port :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
