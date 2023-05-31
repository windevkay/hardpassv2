package main

import (
	"bytes"
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

// execute a page template from cached template sets
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// retrieve template set from cache
	ts, ok := app.templateCache[page]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", page))
		return
	}

	buf := new(bytes.Buffer)
	// write template to buffer
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}