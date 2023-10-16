package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "ERROR \t", log.Ldate|log.Ltime|log.Lshortfile)

  app := &application{
    errorLog: errLog,
    infoLog: infoLog,
  }

	files := http.FileServer(http.Dir("./ui/static/"))
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.view)
	mux.HandleFunc("/snippet/create", create)
	mux.Handle("/static/", http.StripPrefix("/static", files))


  srv := &http.Server {
    Addr: *addr,
    ErrorLog: errLog,
    Handler: mux,
  }
	infoLog.Printf("starting server on: %s", *addr)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}
