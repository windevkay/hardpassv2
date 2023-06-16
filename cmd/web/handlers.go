package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/windevkay/hardpassv2/internal/entities"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	passwords, err := app.passwords.AllPasswords()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Passwords = passwords

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) passwordView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	
	id, err := strconv.Atoi(params.ByName("id"))
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

	data := app.newTemplateData(r)
	data.Password = password

	app.render(w, http.StatusOK, "password.tmpl.html", data)
}

func (app *application) passwordCreatePost(w http.ResponseWriter, r *http.Request) {
	test_app := "test-app"

	id, err := app.passwords.Insert(test_app)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/password/view/%d", id), http.StatusSeeOther)
}


