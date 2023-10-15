package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
  addr := flag.String("addr", ":4000", "HTTP network address")

  flag.Parse()
	files := http.FileServer(http.Dir("./ui/static/"))
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", view)
	mux.HandleFunc("/snippet/create", create)
	mux.Handle("/static/", http.StripPrefix("/static", files))

  log.Printf("starting server on: %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
