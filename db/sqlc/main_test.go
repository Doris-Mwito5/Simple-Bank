package db

import (
	"database/sql"
	"goprojects/simplebank/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load the config", err)
	}
	//db connection using the function Open()
	testDB, err = sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	//the new connection created 
	testQueries = New(testDB)

	//the m.Run executes all the tests an os.Exit ensures the process exits appropriately
	os.Exit(m.Run())
}