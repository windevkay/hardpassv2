package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/windevkay/hardpassv2/internal/entities"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// composed html file paths
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// execute the template
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) passwordCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// set necessary header fields
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	test_app := "test_app"
	id, err := app.passwords.Insert(test_app)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/password/viewOne?id=%d", id), http.StatusSeeOther)
}

func (app *application) passwordViewOne(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	password, err := app.passwords.Get(id)
	if err != nil {
		if errors.Is(err, entities.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/password.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Password: password}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) passwordViewAll(w http.ResponseWriter, r *http.Request) {
	passwords, err := app.passwords.AllPasswords()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, password := range passwords {
		fmt.Fprintf(w, "%+v\n", password)
	}
}
