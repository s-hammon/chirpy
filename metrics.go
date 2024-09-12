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
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	msg := fmt.Sprintf(`<html>

		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
		
		</html>`, a.fserverHits)
	w.Write([]byte(msg))
}

func (a *apiConfig) handleResetMetrics(w http.ResponseWriter, r *http.Request) {
	a.fserverHits = 0
	w.WriteHeader(http.StatusOK)
}
