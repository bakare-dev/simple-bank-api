package util

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/bakare-dev/simple-bank-api/config"
	_ "github.com/lib/pq"
)

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func ConnectDB() *sql.DB {
	dbDriver := "postgres"
	var dbConfig DbConfig

	if config.Settings.Server.Mode == "production" {
		dbConfig = DbConfig{
			Host:     config.Settings.Database.Production.Host,
			Port:     config.Settings.Database.Production.Port,
			User:     config.Settings.Database.Production.User,
			Password: config.Settings.Database.Production.Password,
			Name:     config.Settings.Database.Production.Name,
		}
	} else {
		dbConfig = DbConfig{
			Host:     config.Settings.Database.Development.Host,
			Port:     config.Settings.Database.Development.Port,
			User:     config.Settings.Database.Development.User,
			Password: config.Settings.Database.Development.Password,
			Name:     config.Settings.Database.Development.Name,
		}
	}

	dbSource := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		strconv.Itoa(dbConfig.Port),
		dbConfig.Name,
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
