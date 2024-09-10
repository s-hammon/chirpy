package main

import (
	"log"
	"net/http"
)

func main() {
	const fpathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/assets", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("serving %s on port %s", fpathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
