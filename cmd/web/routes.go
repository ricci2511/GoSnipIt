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

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// register the file server as the handler for all URL pathts
	// starting with "/static/". For matching paths, we strip the
	// "static" prefix before the request reaches the file server.
	r.Handle("/static/", http.StripPrefix("/static", fileServer))


	r.Get("/", app.home)

	// rest routes for snippets
	r.Route("/snippets", func(r chi.Router) {
		r.Get("/create", app.snippetCreateForm)
		r.Post("/", app.snippetCreate)

		r.Route("/{snippetID}", func(r chi.Router) {
			r.Use(app.snippetCtx)
			r.Get("/", app.snippetView)
		})
	})

	return r
}
