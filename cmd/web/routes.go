package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	// simple custom 404 handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// standard middlewares stack to be executed on each request
	r.Use(secureHeaders)
	r.Use(middleware.CleanPath)
	r.Use(app.logRequest)
	r.Use(app.recoverPanic)

	r.Handle("/static/*", staticFileServer("./ui/static/"))

	// all routes run csrf protection and session management middleware
	r.Group(func(r chi.Router) {
		r.Use(noSurf)
		r.Use(app.sessionManager.LoadAndSave)
		r.Use(app.authenticate)

		r.Get("/", app.home)

		// rest routes for user authentication
		r.Route("/user", func(r chi.Router) {
			r.Get("/signup", app.userSignupForm)
			r.Post("/signup", app.userSignup)
			r.Get("/login", app.userLoginForm)
			r.Post("/login", app.userLogin)

			r.With(app.requireAuth).Post("/logout", app.userLogout)
		})

		// rest routes for snippets
		r.Route("/snippets", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(app.requireAuth)
				r.Get("/create", app.snippetCreateForm)
				r.Post("/", app.snippetCreate)
			})

			r.With(app.snippetCtx).Get("/{snippetID}", app.snippetView)
		})
	})

	return r
}

func staticFileServer(staticDir string) http.Handler {
	fs := http.FileServer(http.Dir(staticDir))
	return http.StripPrefix("/static/", fs)
}
