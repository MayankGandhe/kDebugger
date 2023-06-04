package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Print request method and URL
		fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)

		// Print request headers
		fmt.Println("Headers:")
		for name, values := range r.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}

		// Read request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body:", err)

			// Send an error response
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Print request body
		fmt.Println("Body:")
		fmt.Println(string(body))

		// You can add your additional logic here

		// Example response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Fatal(http.ListenAndServe(":5000", nil))
}
