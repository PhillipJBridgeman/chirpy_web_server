/*
 * Author: Phillip Bridgeman
 * Date: 2024-09-03
 * Last Modified Date: 2024-09-03
 * Purpose: This is a simple web server that serves files from the current directory on port 8080.
 */
package main

import (
	"fmt"
	"net/http"
	"sync"
)

type apiConfig struct {
	fileserverHits int
	mu             sync.Mutex
}

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))

	cfg := &apiConfig{}

	// Wrap the fileserver handler with the middleware
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	// Register the handlers with method-based pattern matching
	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("GET /metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /reset", cfg.resetHandler)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	// Start the server
	server.ListenAndServe()
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.mu.Lock()
		cfg.fileserverHits++
		cfg.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.mu.Lock()
	hits := cfg.fileserverHits
	cfg.mu.Unlock()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hits: %d", hits)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.mu.Lock()
	cfg.fileserverHits = 0
	cfg.mu.Unlock()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hits reset to 0"))
}
