package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
	"gosnipit.ricci2511.dev/internal/models"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Fram-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// will recover any panics and send a 500 Internal Server Error response for better UX
// note: this only runs in the goroutine that executed the middleware
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deferred function that will run in the event of a panic
		defer func() {
			if err := recover(); err != nil {
				// this header triggers go's http server to close the connection after a response has been sent
				w.Header().Set("Connection", "close")
				// format the panic value into an error message
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// do not store pages that require authentication in the browser's cache
		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// csrf protection with a custom cookie that is HttpOnly and Secure with path "/"
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

type contextKey string

const (
	contextKeySnippet = contextKey("snippet")
)

func (app *application) snippetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the id from the url param and check if it's valid
		id, err := strconv.Atoi(chi.URLParam(r, "snippetID"))
		if err != nil || id < 1 {
			app.notFound(w)
			return
		}

		// query the db for the snippet
		snippet, err := app.snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
			}

			return
		}

		// add the snippet to the request context and call the next handler
		ctx := context.WithValue(r.Context(), contextKeySnippet, snippet)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
