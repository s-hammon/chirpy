package main

import (
	"log"
	"net/http"

	"github.com/s-hammon/chirpy/internal/database"
)

func main() {
	const fpathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &apiConfig{
		fserverHits: 0,
		DB:          db,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fpathRoot)))))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	mux.HandleFunc("GET /api/reset", cfg.handleResetMetrics)

	mux.HandleFunc("POST /api/chirps", cfg.handleValidateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handleGetChirps)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("serving %s on port %s", fpathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
