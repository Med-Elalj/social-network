package modules

import (
	"database/sql"
	"path/filepath"

	"social-network/app/logs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func SetTables() *sql.DB {
	db, err := sql.Open("sqlite3", "server/db/main.db")
	if err != nil {
		logs.FatalLog.Fatalln("Error opening database:", err)
	}

	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		logs.FatalLog.Fatalln("Failed to set WAL mode:", err)
	}

	// Set the database to use foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		logs.FatalLog.Fatalln("Failed to enable foreign keys:", err)
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		logs.FatalLog.Fatalln("driver instance error:", err)
	}

	absPath, err := filepath.Abs("server/sql/migrations")
	if err != nil {
		logs.FatalLog.Fatalln("unable to resolve migration path:", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+absPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		logs.FatalLog.Fatalln("migration instance error:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logs.FatalLog.Fatalln("migration failed:", err)
	}

	logs.InfoLog.Println("✅ Database migrations applied!")
	return db
}
