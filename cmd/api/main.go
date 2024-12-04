package main

import (
	"log"

	httpServer "github.com/bakare-dev/simple-bank-api/internal/server"

	"github.com/bakare-dev/simple-bank-api/pkg/config"
)

func main() {
	config.LoadConfig()

	httpSvr := httpServer.NewServer()
	if err := httpSvr.Run(); err != nil {
		log.Fatal(err)
	}
}
