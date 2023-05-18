package main

import (
	"errors"
	"fmt"
	"net/http"

	"gosnipit.ricci2511.dev/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
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
	// retrieve the snippet from the context
	ctx := r.Context()
	snippet, ok := ctx.Value(contextKeySnippet).(*models.Snippet)
	if !ok {
		app.serverError(w, errors.New("could not get snippet"))
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display from to create a new snippet..."))
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	id, err := app.snippets.Insert("Dummy title", "Some dummy content, \nIts cool man!", 7)
	if err != nil {
		app.serverError(w, err)
	}

	// redirect user to the page of the newly created snippet
	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}