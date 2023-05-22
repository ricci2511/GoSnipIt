package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"gosnipit.ricci2511.dev/internal/models"
	"gosnipit.ricci2511.dev/ui"
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

	// retrieve a slice of all html pages from the ui embed.FS
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract file name from path
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/partials/*html",
			page,
		}

		// first register any existing template functions and then parse
		// the base template from the embed.FS to the newly created template set
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
