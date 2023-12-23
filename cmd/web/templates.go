package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/HaschwalthB/snippetstash/internal/models"
)

// use Snippets as a dependency for our handlers
// this will make it easier to write our handlers, since we won't have to keep creating a new SnippetModel instance in each handler
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string
}

// make a function time format for snippet
func humanDate(t time.Time) string {
	return t.Local().Format("02 Jan 2006 at 15:04")
}

// create a template.FuncMap to register our custom function
// FuncMap is map[string]any
var function = template.FuncMap{
	"date": humanDate,
}

// make a cache map templateCache
// this handle all our html files
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// find a file path, and get all file path with the extension .html
	// filepath.Glob() returns a slice of strings containing all file paths that match the glob pattern
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	// loop through the pages one-by-one
	// which is the file path like ./ui/html/pages/ .html
	for _, page := range pages {
		// extract the file name (like home.html) from the full file path
		// filepath.Base() returns the last element of the path, which is the file name like home.html, etc
		name := filepath.Base(page)

		// must register our template function before parsing
		tf, err := template.New(name).Funcs(function).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		tf, err = tf.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		tf, err = tf.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// name  containe the file name, that we extracted before in filepath.Base()
		// and we store the template.Template value in the map as the value
		// we  assign tf as value to the key name
		cache[name] = tf
		// name = home.html?
		// tf = template.Template
	}
	return cache, nil
}
