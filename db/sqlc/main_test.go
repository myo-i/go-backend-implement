package db

import (
	"database/sql"
	"go-backend/util"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

// 外部のテストで利用するためにここで宣言
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("connot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("cannot connect to db", err)
	}

	testQueries = New(testDB)

	m.Run()
}
