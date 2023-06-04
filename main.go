package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type DebugHandler struct{}

func (dh DebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	// Convert headers to map[string]string
	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = values[0]
	}

	// Convert headers map to JSON
	headersJSON, err := json.MarshalIndent(headers)
	if err != nil {
		log.Println("Error encoding headers to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Write JSON response
	w.Write(headersJSON)
}

func main() {
	http.Handle("/", DebugHandler{})
	log.Fatal(http.ListenAndServe(":5000", nil))
}
