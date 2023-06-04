package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		// Print request method and URL
		fmt.Printf("Received request: %s %s\n", c.Request.Method, c.Request.URL.Path)

		// Print request headers
		fmt.Println("Headers:")
		for name, values := range c.Request.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}

		// Continue processing the request
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		// Convert headers to map[string]string
		headers := make(map[string]string)
		for name, values := range c.Request.Header {
			headers[name] = values[0]
		}

		// Set Content-Type header
		c.Header("Content-Type", "application/json")

		// Return headers in JSON format
		c.JSON(http.StatusOK, headers)
	})

	log.Fatal(r.Run(":5000"))
}
