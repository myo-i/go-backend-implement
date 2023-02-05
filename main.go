package main

import (
	"database/sql"
	"go-backend/api"
	sqlc "go-backend/db/sqlc"
	"go-backend/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// "."現在のフォルダ
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("cannot connect to db", err)
	}

	store := sqlc.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create new server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
