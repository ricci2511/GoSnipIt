package main

import (
	"html/template"
	"path/filepath"
	"time"

	"gosnipit.ricci2511.dev/internal/models"
)

// holds any dynamic data that we want to pass to our HTML templates
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string // holds flash messages
	IsAuthenticated bool
	CSRFToken       string
}

// formats dates in a human-readable format
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// global variable to hold the functions that we want to make available in our templates
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// initialiazes a map to hold the template set of all pages with the page file name as the key
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract file name from path
		name := filepath.Base(page)

		// first register the template functions and then add
		// the base template to the newly created template set
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
