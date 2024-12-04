package db

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/bakare-dev/simple-bank-api/config"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	config.LoadConfig()

	dbDriver := "postgres"
	dbSource := "postgresql://" +
		config.Settings.Database.User + ":" +
		config.Settings.Database.Password + "@" +
		config.Settings.Database.Host + ":" +
		strconv.Itoa(config.Settings.Database.Port) + "/" +
		config.Settings.Database.Name + "?sslmode=disable"

	var err error

	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
