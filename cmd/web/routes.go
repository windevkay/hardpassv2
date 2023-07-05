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
	withSession := alice.New(app.sessionManager.LoadAndSave, noSurf)
	// public routes
	router.Handler(http.MethodGet, "/user/signup", withSession.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", withSession.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", withSession.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", withSession.ThenFunc(app.userLoginPost))

	// auth guard middleware
	withAuth := withSession.Append(app.requireAuthentication)
	// private routes
	router.Handler(http.MethodPost, "/password/create", withAuth.ThenFunc(app.passwordCreatePost))
	router.Handler(http.MethodPost, "/user/logout", withAuth.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/", withAuth.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/password/view/:id", withAuth.ThenFunc(app.passwordView))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
