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
	ShipmentURL string `envconfig:"SHIPMENT_URL" required:"true"`
	Port       string `envconfig:"PORT" default:"8084"`
}

// healthHandler responds with HTTP 200 for health checks
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy")) // Add a simple body response for better debugging
}

// corsMiddleware adds CORS headers to the response
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust the allowed origin as needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight (OPTIONS) requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
    var config AppConfig
    if err := envconfig.Process("", &config); err != nil {
        log.Fatalf("Failed to parse environment variables: %v", err)
    }

    server, err := NewGraphQLServer(
        config.AccountURL,
        config.ShopifyURL,
        config.ShipmentURL,
    )
    if err != nil {
        log.Fatalf("Failed to create GraphQL server: %v", err)
    }

    http.Handle("/graphql", corsMiddleware(handler.GraphQL(server.ToNewExecutableSchema())))
    http.Handle("/playground", corsMiddleware(playground.Handler("GraphQL Playground", "/graphql")))
    http.Handle("/health", corsMiddleware(http.HandlerFunc(healthHandler)))

    log.Printf("Starting server on port %s...", config.Port)
    if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
