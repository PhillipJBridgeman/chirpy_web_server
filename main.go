package main

import (
	"net/http"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(".")) // Serve files from the current directory

	mux.Handle("/healthz", http.HandlerFunc(healthzHandler))
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	// Start the server
	server.ListenAndServe()
}
