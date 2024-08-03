package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(".")) // Serve files from the current directory

	mux.Handle("/", fileServer)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	// Start the server
	server.ListenAndServe()
}
