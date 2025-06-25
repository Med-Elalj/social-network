package modules

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"social-network/server/logs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func SetDb(db *sql.DB) {
	DB = db
}

func SetTables() *sql.DB {
	db, err := sql.Open("sqlite3", "server/db/main.db")
	if err != nil {
		logs.Fatalf("Error opening database: %v", err)
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatal("driver instance error:", err)
	}

	absPath, err := filepath.Abs("server/sql/migrations")
	if err != nil {
		log.Fatal("unable to resolve migration path:", err)
	}
	fmt.Println("Migration path:", absPath)
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+absPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal("migration instance error:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration failed:", err)
	}

	logs.Println("✅ Database migrations applied!")
	DB = db
	return db
}
