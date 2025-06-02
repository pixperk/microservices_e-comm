package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pixperk/microservices_e-comm/catalog"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
		}
		return
	})

	defer r.Close()
	log.Println("Starting catalog service on port 8080...")
	s := catalog.NewService(r)
	if err := catalog.ListenGRPC(s, 8080); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

}
