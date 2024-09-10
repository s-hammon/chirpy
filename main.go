package main

import (
	"log"
	"net/http"
)

func main() {
	const fpathRoot = "."
	const port = "8080"

	cfg := &apiConfig{
		fserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fpathRoot)))))
	mux.HandleFunc("/healthz", handleReadiness)
	mux.HandleFunc("/metrics", cfg.handleMetrics)
	mux.HandleFunc("/reset", cfg.handleResetMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("serving %s on port %s", fpathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
