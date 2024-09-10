package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fserverHits int
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fserverHits++
		next.ServeHTTP(w, r)
	})
}

func (a *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	msg := fmt.Sprintf("Hits: %d\n", a.fserverHits)
	w.Write([]byte(msg))
}

func (a *apiConfig) handleResetMetrics(w http.ResponseWriter, r *http.Request) {
	a.fserverHits = 0
	w.WriteHeader(http.StatusOK)
}
