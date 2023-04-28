package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)


func main() {
	// command line flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// levelled logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/password/create", passwordCreate)
	mux.HandleFunc("/password/view", passwordView)

	infoLog.Printf("Hardpass is starting on port %s 🚀", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}
