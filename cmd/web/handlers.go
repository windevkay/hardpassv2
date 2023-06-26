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
	App                 string `form:"app"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	Name                string `form:"name"`
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

	app.sessionManager.Put(r.Context(), "flash", "Password successfully generated!")

	http.Redirect(w, r, fmt.Sprintf("/password/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
