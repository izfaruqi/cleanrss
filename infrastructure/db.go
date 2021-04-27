package infrastructure

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", "./db.sqlite3")
	if err != nil {
		return nil, err
	}
	return db, nil
}
