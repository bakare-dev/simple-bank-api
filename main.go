package main

import (
	"fmt"

	"github.com/bakare-dev/simple-bank-api/config"
)

func main() {
	config.LoadConfig()

	fmt.Printf("Server running at %s:%d\n", config.Settings.Server.Host, config.Settings.Server.Port)
}
