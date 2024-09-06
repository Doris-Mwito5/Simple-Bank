package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	_ "github.com/lib/pq"
)

//constants to open a connection to your PostgreSQL database.
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5434/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	//db connection using the function Open()
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	//the new connection created 
	testQueries = New(conn)

	//the m.Run executes all the tests an os.Exit ensures the process exits appropriately
	os.Exit(m.Run())
}