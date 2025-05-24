package db

import (
	"database/sql"
)

var DB *sql.DB

func SetDb(db *sql.DB) {
	DB = db
}
