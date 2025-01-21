package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/conneroisu/groq-go/extensions/e2b"
)

func main() {
	// Initialize the middleware
	middleware, err := NewMiddleware("postgres://user:password@localhost:5432/dbname")
	if err != nil {
		log.Fatalf("Failed to initialize middleware: %v", err)
	}
	defer middleware.Close()

	// Initialize the e2b client
	e2bClient, err := e2b.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize e2b client: %v", err)
	}

	// Set up the router
	router := httprouter.New()

	// Define routes and handlers
	router.POST("/query", middleware.HandleQuery)
	router.POST("/subtract", middleware.HandleSubtract)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
