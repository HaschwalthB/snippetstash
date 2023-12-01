package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// initialize new httprouter
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	// handle static file
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.view)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetPost)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetForm)

	standart := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	return standart.Then(router)
}
