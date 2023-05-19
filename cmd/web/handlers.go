package main

import (
	"errors"
	"fmt"
	"net/http"

	"gosnipit.ricci2511.dev/internal/models"
	"gosnipit.ricci2511.dev/internal/validator"
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
	// initialize a basic templateData and render the snipet create form
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 7,
	}
	app.render(w, http.StatusOK, "create.html", data)
}

// add struct tags to the fields to tell the form decoder how to map the form data to the struct
type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // tell decoder to ignore this field
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// parses form data into r.PostForm map
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form snippetCreateForm

	// decode and fill the form struct with the relevant form data
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// form validation
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be longer than 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must be either 1, 7 or 365")

	if !form.Valid() {
		// create new templateData with the populated errors and render the form again
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
	}

	// redirect user to the page of the newly created snippet
	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}
