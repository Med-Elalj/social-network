package main

import (
	"database/sql"
	"os"

	"social-network/server/logs"
)

var db *sql.DB

// func init() {
// 	// // Open logs. files
// 	// logs.File, err := os.OpenFile("./Serverlogs.s/app.logs.", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
// 	// if err != nil {
// 	// 	logs..Fatalf("Failed to open app logs. file: %v", err)
// 	// }

// 	// errorFile, err := os.OpenFile("./Serverlogs.s/error.logs.", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
// 	// if err != nil {
// 	// 	logs..Fatalf("Failed to open error logs. file: %v", err)
// 	// }

// 	// // Setup logs.gers
// 	// helpers.Infologs. = logs..New(logs.File, "INFO: ", logs..Ldate|logs..Ltime|logs..Lshortfile)
// 	// helpers.Errorlogs. = logs..New(errorFile, "ERROR: ", logs..Ldate|logs..Ltime|logs..Lshortfile)

// 	// // Set default logs. output for Fatal logs.s
// 	// logs..SetOutput(errorFile)
// 	db = SetTables()
// }

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
	db = dB
	return dB
}
