package main

import (
	"log"
	"time"

	"github.com/Shridhar2104/logilo/account"

	"github.com/tinrab/retry"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_ACCOUNT_URL"`
}

func main() {

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to process environment variables: %v", err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {

		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return err
		}
		return nil
	})
	defer r.Close()
	log.Println("server starting on port 8081 ...")
	  

	s := account.NewAccountService(r)
	log.Fatal(account.NewGRPCServer(s, 8081))

	
}
