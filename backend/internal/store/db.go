package store

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type SqlDb struct {
	db *sql.DB
}

func NewSqlDB(url string) *SqlDb {
	db, err := sql.Open("sqlite", url)

	if err != nil {
		log.Fatal("Failed open db connection", err)
	}

	return &SqlDb{
		db: db,
	}
}
