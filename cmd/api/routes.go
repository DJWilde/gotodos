package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	// User routes
	router.HandlerFunc(http.MethodPost, "/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/login", app.createTokenHandler)

	return router
}