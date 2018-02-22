package main

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/vardius/blockchain/pkg/client"
	"github.com/vardius/blockchain/pkg/wallet"
)

type config struct {
	Host       string `env:"HOST" envDefault:"localhost"`
	Port       int    `env:"PORT" envDefault:"80"`
	ServerHost string `env:"SERVER_HOST" envDefault:"localhost"`
	ServerPort int    `env:"SERVER_PORT" envDefault:"3000"`
}

func main() {
	err := godotenv.Load("./cmd/wallet/.env")
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{}
	env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	c := client.New(cfg.ServerHost, cfg.ServerPort)
	w := wallet.NewHTTP(c)

	w.Run(cfg.Host, cfg.Port)
}
