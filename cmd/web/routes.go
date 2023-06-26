package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	// not found can wrap our helper function
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	// session middleware
	withSession := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", withSession.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/password/view/:id", withSession.ThenFunc(app.passwordView))
	router.Handler(http.MethodPost, "/password/create", withSession.ThenFunc(app.passwordCreatePost))
	router.Handler(http.MethodGet, "/user/signup", withSession.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", withSession.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", withSession.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", withSession.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", withSession.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
