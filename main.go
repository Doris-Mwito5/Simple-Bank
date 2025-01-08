package main

import (
	"database/sql"
	"goprojects/simplebank/api"
	db "goprojects/simplebank/db/sqlc"
	"goprojects/simplebank/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config, err")
	}
	//db connection using the function Open()
	conn, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.SERVER_ADDRESS)
	if err != nil {
		log.Fatal("could not start the server", err)
	}


}