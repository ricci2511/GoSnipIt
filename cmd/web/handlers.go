package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// since "/" ends with a trailing slash, it uses subtree matching
	// meaning that any URL paths that start with "/" will match,
	// therefore we check if the current request URL path exactly matches "/"
    if r.URL.Path != "/" {
		app.notFound(w)
        return
    }

	// important: the file containing the base template must be the first file in the slice
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	// read the html files into a template set
	// files is passed as a variadic parameter
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// write the content of the base template as the response body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
		app.notFound(w)
        return
    }

    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
        return
    }

    w.Write([]byte("Create a new snippet..."))
}