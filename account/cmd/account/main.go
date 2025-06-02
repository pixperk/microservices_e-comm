package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pixperk/microservices_e-comm/account"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconnfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
		}
		return
	})

	defer r.Close()
	s := account.NewService(r)
	log.Println("Starting account service on port 8080...")
	if err := account.ListenGRPC(s, 8080); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

}
