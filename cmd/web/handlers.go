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

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)
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

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)
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