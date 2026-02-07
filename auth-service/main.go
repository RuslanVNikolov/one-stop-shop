package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Simple handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from One Stop Shop Auth Service!")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	// Start server
	port := ":8001"
	fmt.Printf("Auth Service starting on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
