package db

import (
	"database/sql"
	"log"
)

const TestDBConn = "postgres://postgres:postgrespw@localhost:32768/gochat?sslmode=disable"

func OpenTestDB() (*sql.DB, *sql.Tx, error) {

	db, err := sql.Open("postgres", TestDBConn)
	if err != nil {
		log.Fatalf("Failed to open test DB connection: %s", err)
		return nil, nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %s", err)
		return nil, nil, err
	}

	return db, tx, nil
}

func CloseTestDB(tx *sql.Tx, db *sql.DB) {
	tx.Rollback()
	db.Close()
}
