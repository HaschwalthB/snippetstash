package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/HaschwalthB/snippetstash/internal/models"
	"github.com/julienschmidt/httprouter"
)

// use the application struct to hold the application-wide dependencies for the web application
// for now, it'll only contain fields for the two custom loggers, and database model.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// call templateData method that return templateData struct
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) view(w http.ResponseWriter, r *http.Request) {
	// retrieve router parameters from context
	// which is the id of the snippet (/snippet/view/:id)
	// convert the id parameter from a string to an integer
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.html", data)

}

// snippetNew display the form for creating a new snippet
func (app *application) snippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "create.html", data)
}

// snippetPost handler and a event handler for the form submission
func (app *application) snippetPost(w http.ResponseWriter, r *http.Request) {
	// parse the body request
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// get the form data after parse, this return string
	// so for context of expires we need to convert it to int
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}
	// make a validation for the form
	falidErrors := make(map[string]string)
	// check title for empty string and long character
	if strings.TrimSpace(title) == "" {
		falidErrors["title"] = "the field cannot blank my friend"
	} else if utf8.RuneCountInString(title) > 50 {
		falidErrors["title"] = "to much!!! "
	}
	if strings.TrimSpace(content) == "" {
		falidErrors["title"] = "cmon brohh!!! are you st**id or what. do you wanna make a snippet but you dont filled this up?. get a docter!"
		return
	}
	if expires != 1 && expires != 7 && expires != 365 {
		falidErrors["expires"] = "just choose one"
		return
	}

	if len(falidErrors) > 0 {
		fmt.Fprint(w, falidErrors)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}
	// make a clean path to the snippetview page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
