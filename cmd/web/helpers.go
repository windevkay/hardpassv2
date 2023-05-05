package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError logs the error and sends a generic 500 Internal Server Error response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	// stack trace
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// log the error
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a specific status code and description to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	// write status code
	http.Error(w, http.StatusText(status), status)
}

// notFound sends a 404 Not Found response to the user
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
