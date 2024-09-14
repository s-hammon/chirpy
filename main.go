package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/s-hammon/chirpy/internal/database"
)

const dbPath = "database.json"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	const fpathRoot = "."
	const port = "8080"

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		os.Remove(dbPath)
	}

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &apiConfig{
		fserverHits: 0,
		DB:          db,
		jwtSecret:   os.Getenv("JWT_SECRET"),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fpathRoot)))))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	mux.HandleFunc("GET /api/reset", cfg.handleResetMetrics)

	mux.HandleFunc("POST /api/chirps", cfg.handleCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handleGetChirpByID)

	mux.HandleFunc("POST /api/users", cfg.handleNewUser)
	mux.HandleFunc("PUT /api/users", cfg.handleUpdateUser)
	mux.HandleFunc("POST /api/login", cfg.handleLogin)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("serving %s on port %s", fpathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
