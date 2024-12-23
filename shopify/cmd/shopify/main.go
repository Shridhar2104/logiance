package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/Shridhar2104/logilo/shopify"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_SHOPIFY_URL"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to process environment variables: %v", err)
	}

	var db *sql.DB
	var err error

	// Retry database connection
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		db, err = sql.Open("postgres", cfg.DatabaseURL)
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return err
		}

		// Test the connection
		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping database: %v", err)
			return err
		}

		return nil
	})
	defer db.Close()

	// Initialize repository with the database connection
	r := shopify.NewPostgresRepository(db)

	log.Println("server starting on port 8080 ...")

	ApiKey := "67f10611ac39283d047c7cc4c8e04954"
	ApiSecret := "63e9e494ff13bddd03cb4d742baa10f0"
	RedirectUrl := "http://localhost:3000/storeorders"

	s := shopify.NewShopifyService(ApiKey, ApiSecret, RedirectUrl, r)
	log.Fatal(shopify.NewGRPCServer(s, 8080))
}