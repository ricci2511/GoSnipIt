package main

import "net/http"

func (app *application) routes() *http.ServeMux {
    mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// register the file server as the handler for all URL pathts
	// starting with "/static/". For matching paths, we strip the
	// "static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// application routes
    mux.HandleFunc("/", app.home)
    mux.HandleFunc("/snippet/view", app.snippetView)
    mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}