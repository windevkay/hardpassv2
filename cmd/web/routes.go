package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/password/create", app.passwordCreate)
	mux.HandleFunc("/password/view", app.passwordView)

	return mux
}
