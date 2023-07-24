package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/windevkay/hardpassv2/internal/azure"
	"github.com/windevkay/hardpassv2/internal/entities"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	passwords      entities.PasswordEntityInterface
	users          entities.UserEntityInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// authenticate and setup azure client
	client, err := azure.SetupClient()
	if err != nil {
		log.Fatal(err)
	}

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

	formDecoder := form.NewDecoder()

	// session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = true // only send cookies over HTTPS

	// initialize application struct
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		passwords:      &entities.PasswordEntity{DB: db, AzureClient: client},
		users:          &entities.UserEntity{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// custom tls config - limit elliptic curves in tls handshake
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// override some server defaults
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,      // close connections after 1 minute of inactivity
		ReadTimeout:  5 * time.Second,  // allow 5 seconds to read request headers
		WriteTimeout: 10 * time.Second, // allow 10 seconds to write response
	}

	infoLog.Printf("Hardpass is starting on port %s 🚀", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
