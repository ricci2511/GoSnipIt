package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// writes an error message and stack trace to the errorLog, then
// sends a generic 500 Internal Server Error response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// sends a specific status code and corresponding message to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// common helper for 404 Not Found responses, convenience wrapper around clientError()
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// first check if the template exists in the cache
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// init buffer to hold template output
	buf := new(bytes.Buffer)

	// write the template to the buffer to catch any errors before
	// writing directly to the http.ResponseWriter
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	// finally write the contents of the buffer to the http.ResponseWriter
	buf.WriteTo(w)
}

// helper to decode form data into a struct (target being the struct to decode into)
func (app *application) decodePostForm(r *http.Request, target any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(target, r.PostForm)
	if err != nil {
		// if the error is an InvalidDecoderError, we probably passed in an invalid struct
		// so errors.As() is used to check and panic if that's the case
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// other types of errors are likely caused by the client
		return err
	}

	return nil
}
