package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://<user>:<password>@localhost:5432/simplebank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	
	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Cannot create db: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
