package main

import (
	"log"
	"net/http"
	

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_URL" required:"true"`
	ShopifyURL string `envconfig:"SHOPIFY_URL" required:"true"`
	Port       string `envconfig:"PORT" default:"8084"` 
}

// healthHandler responds with HTTP 200 for health checks
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy")) // Add a simple body response for better debugging
}


func main() {
	// Load configuration from environment variables
	var config AppConfig
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}

	// Create a new GraphQL server
	server, err := NewGraphQLServer(config.AccountURL, config.ShopifyURL)
	if err != nil {
		log.Fatalf("Failed to create GraphQL server: %v", err)
	}

	// Set up HTTP routes
	http.Handle("/graphql", handler.GraphQL(server.ToNewExecutableSchema()))
	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
	http.Handle("/health", http.HandlerFunc(healthHandler))

	// Start the server
	log.Printf("Starting server on port %s...", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
