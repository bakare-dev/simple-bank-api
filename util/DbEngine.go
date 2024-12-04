package util

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/bakare-dev/simple-bank-api/config"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	dbDriver := "postgres"
	dbSource := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		config.Settings.Database.User,
		config.Settings.Database.Password,
		config.Settings.Database.Host,
		strconv.Itoa(config.Settings.Database.Port),
		config.Settings.Database.Name,
	)

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	if err = conn.Ping(); err != nil {
		log.Fatalf("cannot verify db connection: %v", err)
	}

	return conn
}
