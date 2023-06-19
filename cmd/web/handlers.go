package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/windevkay/hardpassv2/internal/entities"
	"github.com/windevkay/hardpassv2/internal/validator"

	"github.com/julienschmidt/httprouter"
)

type passwordCreateForm struct {
	App string `form:"app"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	passwords, err := app.passwords.AllPasswords()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Passwords = passwords
	data.Form = passwordCreateForm{}

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
	var form passwordCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.App), "app", "App name cannot be blank")
	form.CheckField(validator.MaxChars(form.App, 50), "app", "App name cannot be longer than 50 characters")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "home.tmpl.html", data)
		return
	}

	id, err := app.passwords.Insert(form.App)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/password/view/%d", id), http.StatusSeeOther)
}


