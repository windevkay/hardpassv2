package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// command line flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// levelled logs - concurrency safe
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// initialize application struct
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// override some server defaults
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Hardpass is starting on port %s 🚀", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
