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
		config.Settings.Database.Test.User + ":" +
		config.Settings.Database.Test.Password + "@" +
		config.Settings.Database.Test.Host + ":" +
		strconv.Itoa(config.Settings.Database.Test.Port) + "/" +
		config.Settings.Database.Test.Name + "?sslmode=disable"

	var err error

	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
