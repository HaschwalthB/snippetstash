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

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/schema"

	"github.com/HaschwalthB/snippetstash/internal/models"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModelDB
	users          *models.UserModelDB
	templateCache  map[string]*template.Template
	schema         *schema.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":9000", "HTTP network address")

	dsn := flag.String("dsn", "web:komeng@/snippetbox?parseTime=true", "MySQL")
	// parse the command-line flag before using it
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR \t", log.Ldate|log.Ltime|log.Lshortfile)

	// create a connection pool to the database openDB function returns a sql.DB connection pool
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// close db connection before main() exits
	defer db.Close()

	// initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// initialize a new schema decoder
	decoder := schema.NewDecoder()

	// initialize a new session manager
	// we use our database as a session manage
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	// initialize a new instance of application containing the dependencies
	// snippets its a pointer to a SnippetModel struct, which holds a sql.DB connection pool
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModelDB{DB: db},
		users:          &models.UserModelDB{DB: db},
		templateCache:  templateCache,
		schema:         decoder,
		sessionManager: sessionManager,
	}

	// change the default setting of tls config3
	tlsConfig := tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: &tlsConfig,

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// openDB() wraps sql.Open() and returns a sql.DB connection pool
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// check the connection by pinging the database
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
