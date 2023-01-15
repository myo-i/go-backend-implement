package db

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5433/test_bank?sslmode=disable"
)

var testQueries *Queries

// 外部のテストで利用するためにここで宣言
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalln("cannot connect to db", err)
	}

	testQueries = New(testDB)

	m.Run()
}
