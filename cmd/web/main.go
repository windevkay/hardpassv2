package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/windevkay/hardpassv2/internal/entities"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	passwords     *entities.PasswordEntity
	templateCache map[string]*template.Template
}

func main() {
	// defaults from env
	defaultPort := os.Getenv("PORT")
	defaultDSN := os.Getenv("DSN")

	// command line flags
	addr := flag.String("addr", defaultPort, "HTTP network address")
	dsn := flag.String("dsn", defaultDSN, "MySQL data source name")
	flag.Parse()

	// levelled logs - concurrency safe
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// database connection pool
	db, err := getDBConnectionPool(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// initialize application struct
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		passwords:     &entities.PasswordEntity{DB: db},
		templateCache: templateCache,
	}

	// override some server defaults
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Hardpass is starting on port %s 🚀", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func getDBConnectionPool(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// set max open connections
	db.SetMaxOpenConns(25)
	// set max idle connections
	db.SetMaxIdleConns(25)
	// set max connection lifetime
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}
