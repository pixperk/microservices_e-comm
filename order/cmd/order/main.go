package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pixperk/microservices_e-comm/order"

	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL"`
	Port        int    `envconfig:"PORT"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to process environment variables: %v", err)
	}

	var r order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewPostgresRepository(cfg.DatabaseURL)
		return
	})
	defer r.Close()

	log.Println("Order service started on port 8080...")
	s := order.NewService(r)
	if err := order.ListenGRPC(s, cfg.AccountURL, cfg.CatalogURL, cfg.Port); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

}
