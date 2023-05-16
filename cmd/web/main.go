package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// struct to hold the application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	// cli flags
	// fall back to 4000 if no port is provided
	addr := flag.String("addr", ":4000", "HTTP network adress")   

	// this reads the the value of the flags and assigns them to their variables
	// important: must be called before using the variables, e.g addr
	flag.Parse()

	// info logger that includes date and time
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// error logger that includes date, time, source file name and line number
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// init new application struct
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	// init new http.Server struct with the port, custom error logger and routes
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

    infoLog.Printf("Starting server on %s", *addr)
    err := srv.ListenAndServe()
    errorLog.Fatal(err)
}