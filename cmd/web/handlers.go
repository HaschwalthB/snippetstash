package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/HaschwalthB/snippetstash/internal/models"
	"github.com/HaschwalthB/snippetstash/internal/validator"
)

// make a object for validation form
type snippetCreateForm struct {
	Title               string `schema:"title"`
	Content             string `schema:"content"`
	Expires             int    `schema:"expires"`
	validator.Validator `schema:"-"`
}

// use the application struct to hold the application-wide dependencies for the web application
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

	data.Form = &snippetCreateForm{
		Expires: 365,
	}
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

	var schema snippetCreateForm
	// get the form data after parse from body request, this return string
	// so for context of expires we need to convert it to int
	err = app.schema.Decode(&schema, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	schema.Checkvield(validator.NotBlank(schema.Title), "title", "field required")
	schema.Checkvield(validator.MaxChars(schema.Title, 100), "title", "field must be less than 50 characters")
	schema.Checkvield(validator.NotBlank(schema.Content), "content", "field required")
	schema.Checkvield(validator.PermittedInt(schema.Expires, 1, 7, 365), "Expires", "invalid value")

	if !schema.Valid() {
		data := app.newTemplateData(r)
		data.Form = &schema
		app.render(w, http.StatusOK, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(schema.Title, schema.Content, schema.Expires)
	if err != nil {
		app.serverError(w, err)
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet Create susccesfully!")
	// make a clean path to the snippetview page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
