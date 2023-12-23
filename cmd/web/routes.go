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

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.view))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetPost))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetForm))

	standart := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	return standart.Then(router)
}
