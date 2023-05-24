package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"gosnipit.ricci2511.dev/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// struct to hold the application-wide dependencies
type application struct {
	debug          bool
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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
	debug := flag.Bool("debug", false, "Enable debug mode")

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

	formDecoder := form.NewDecoder()

	// init session manager, use mysql as session store, expires after 12 hours
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// set secure flag on session cookie so it's only sent over HTTPS
	sessionManager.Cookie.Secure = true

	// init new application struct
	app := &application{
		debug:          *debug,
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// restrict elliptic curves to X25519 and P256 which have assembly implementations,
	// therefore they're less cpu intensive than other curves
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// init new http.Server struct with the host, custom error logger and routes
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
