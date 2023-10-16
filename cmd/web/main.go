package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)

	errLog := log.New(os.Stdout, "ERROR \t", log.Ldate|log.Ltime|log.Lshortfile)
	files := http.FileServer(http.Dir("./ui/static/"))
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", view)
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
