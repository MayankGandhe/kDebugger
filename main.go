package main

import (
	"log"
	"net/http"
)

type DebugHandler struct{}

func (dh DebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	// Log all headers
	for name, headers := range r.Header {
		for _, h := range headers {
			log.Printf("Header: %s: %s", name, h)
		}
	}

	// Add your additional logic here

	// Example response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.Handle("/", DebugHandler{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
