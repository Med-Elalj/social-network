package modules

import (
	"database/sql"
	"os"

	"social-network/server/logs"
)

var DB *sql.DB

func SetDb(db *sql.DB) {
	DB = db
}

func SetTables() *sql.DB {
	dB, err := sql.Open("sqlite3", "server/db/main.db")
	if err != nil {
		logs.Fatalf("Error opening database: %v", err)
	}

	sqlContent, err := os.ReadFile("./server/sql/db.sql")
	if err != nil {
		logs.Fatalf("Error reading SQL file: %v", err)
	}

	_, err = dB.Exec(string(sqlContent))
	if err != nil {
		logs.Fatalf("Error executing SQL: %v", err)
	}
	logs.Println("Database successfully created!")
	DB = dB
	return dB
}
