package db

import (
	"gochatv1/config"

	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	db, err := sql.Open(cfg.DB, cfg.DBConn)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
