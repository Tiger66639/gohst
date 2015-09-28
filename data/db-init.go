package data

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect(connection string) {
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}
}
