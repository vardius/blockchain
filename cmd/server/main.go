package main

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/vardius/blockchain/pkg/server"
)

type config struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"3000"`
}

func main() {
	err := godotenv.Load("./cmd/server/.env")
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{}
	env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Run(cfg.Host, cfg.Port))
}
