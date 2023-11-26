package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.view)
	mux.HandleFunc("/snippet/create", app.create)

	standart := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	return standart.Then(mux)
}
