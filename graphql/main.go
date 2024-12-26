package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
    AccountURL  string `envconfig:"ACCOUNT_URL" required:"true"`
    ShopifyURL  string `envconfig:"SHOPIFY_URL" required:"true"`
    PaymentURL  string `envconfig:"PAYMENT_URL" required:"true"`
    ShipmentURL string `envconfig:"SHIPMENT_URL" required:"true"`
    Port        string `envconfig:"PORT" default:"8084"`
}

// healthHandler responds with HTTP 200 for health checks
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Healthy"))
}

// corsMiddleware adds CORS headers to the response
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
    // Initialize logger
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    // Print current working directory for debugging
    pwd, _ := os.Getwd()
    log.Printf("Current working directory: %s", pwd)

    // Load .env.development file
    envPath := "/app/.env.development"
    log.Printf("Attempting to load env file from: %s", envPath)
    
    if err := godotenv.Load(envPath); err != nil {
        log.Fatalf("Error loading .env.development file: %v", err)
    }

    // Process environment variables
    var config AppConfig
    if err := envconfig.Process("", &config); err != nil {
        log.Fatalf("Failed to parse environment variables: %v", err)
    }

    // Create a new GraphQL server
    server, err := NewGraphQLServer(config.AccountURL, config.ShopifyURL, config.ShipmentURL, config.PaymentURL)
    if err != nil {
        log.Fatalf("Failed to create GraphQL server: %v", err)
    }

    // Set up routes
    http.Handle("/graphql", corsMiddleware(handler.GraphQL(server.ToNewExecutableSchema())))
    http.Handle("/playground", corsMiddleware(playground.Handler("GraphQL Playground", "/graphql")))
    http.Handle("/health", corsMiddleware(http.HandlerFunc(healthHandler)))

    log.Printf("Starting server on port %s...", config.Port)
    if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}