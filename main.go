package main

import (
	"log"
	"strconv"

	"github.com/bakare-dev/simple-bank-api/api"
	"github.com/bakare-dev/simple-bank-api/config"
	db "github.com/bakare-dev/simple-bank-api/db/sqlc"
	"github.com/bakare-dev/simple-bank-api/util"
)

func main() {
	config.LoadConfig()

	serverAddress := config.Settings.Server.Host + ":" + strconv.Itoa(config.Settings.Server.Port)

	conn := util.ConnectDB()

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(serverAddress); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
