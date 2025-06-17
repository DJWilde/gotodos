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

	// Todo routes
	router.HandlerFunc(http.MethodGet, "/user/:userId/todos", app.requireAuthenticatedUser(app.getTodosByUserIdHandler))
	router.HandlerFunc(http.MethodGet, "/todos/:id", app.requireAuthenticatedUser(app.getTodoByIdHandler))
	router.HandlerFunc(http.MethodPost, "/todos", app.requireAuthenticatedUser(app.createTodoHandler))
	router.HandlerFunc(http.MethodPatch, "/todos/:id", app.requireAuthenticatedUser(app.updateTodoHandler))
	router.HandlerFunc(http.MethodDelete, "/todos/:id", app.requireAuthenticatedUser(app.deleteTodoHandler))

	return app.authenticate(router)
}