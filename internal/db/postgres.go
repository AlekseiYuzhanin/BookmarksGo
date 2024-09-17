package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		content TEXT NOT NULL
	);`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}
