package main

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func init() {
	// // Open log files
	// logFile, err := os.OpenFile("./ServerLogs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	// if err != nil {
	// 	log.Fatalf("Failed to open app log file: %v", err)
	// }

	// errorFile, err := os.OpenFile("./ServerLogs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	// if err != nil {
	// 	log.Fatalf("Failed to open error log file: %v", err)
	// }

	// // Setup loggers
	// helpers.InfoLog = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	// helpers.ErrorLog = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// // Set default log output for Fatal logs
	// log.SetOutput(errorFile)
	db = SetTables()
}

func SetTables() *sql.DB {
	db, err := sql.Open("sqlite3", "./DataBase/SN.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	sqlContent, err := os.ReadFile("./DataBase/db.DB.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}

	_, err = db.Exec(string(sqlContent))
	if err != nil {
		log.Fatalf("Error executing SQL: %v", err)
	}
	log.Println("Database successfully created!")
	return db
}
