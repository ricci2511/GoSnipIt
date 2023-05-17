package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"gosnipit.ricci2511.dev/internal/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// struct to hold the application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// info logger that includes date and time
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// error logger that includes date, time, source file name and line number
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// load env variables from .env file
	env, err := getDotEnvVars()
	if err != nil {
		errorLog.Fatal(err)
	}

	// cli flags
	// fall back to 4000 if no host is provided
	addr := flag.String("addr", ":4000", "HTTP network adress")   

	// default dsn connection string
	dsnStr := fmt.Sprintf("%v:%v@/gosnipit?parseTime=true", env["MYSQL_USER"], env["MYSQL_PASSWORD"])
	dsn := flag.String("dsn", dsnStr, "MySQL database connection string")

	// this reads the the value of the flags and assigns them to their variables
	// important: must be called before using the variables, e.g addr
	flag.Parse()

	db, err := openDb(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// defer closing the db connection pool until the main() function has finished
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// init new application struct
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// init new http.Server struct with the host, custom error logger and routes
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

    infoLog.Printf("Starting server on %s", *addr)
    err = srv.ListenAndServe()
    errorLog.Fatal(err)
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	} else {
		log.Println("Connected to database")
	}
	

	return db, nil
}

func getDotEnvVars() (map[string]string, error) {
	// load .env file into a map
	envMap, err := godotenv.Read()
	if err != nil {
		return nil, err
	}

	return envMap, nil
}
