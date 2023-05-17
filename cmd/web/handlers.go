package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gosnipit.ricci2511.dev/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// since "/" ends with a trailing slash, it uses subtree matching
	// meaning that any URL paths that start with "/" will match,
	// therefore we check if the current request URL path exactly matches "/"
    if r.URL.Path != "/" {
		app.notFound(w)
        return
    }

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	/*
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
	*/
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

    fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
        return
    }

	id, err := app.snippets.Insert("Dummy title", "Some dummy content, \nIts cool man!", 7)
	if err != nil {
		app.serverError(w, err)
	}

	// redirect user to the page of the newly created snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}