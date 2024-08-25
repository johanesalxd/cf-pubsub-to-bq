package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/johanesalxd/cf-pubsub-to-bq"
)

// main is the entry point of the application.
// It sets up the server port and starts the Functions Framework.
func main() {
	// Set default port to 8080
	port := "8080"

	// Override port if PORT environment variable is set
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// Start the Functions Framework
	// This will listen for incoming HTTP requests and route them to the appropriate function
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
