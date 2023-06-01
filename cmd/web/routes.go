package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/password/create", app.passwordCreate)
	mux.HandleFunc("/password/viewOne", app.passwordViewOne)

	return app.logRequest(secureHeaders(mux))
}
