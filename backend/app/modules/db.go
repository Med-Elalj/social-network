package modules

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"social-network/server/logs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func SetTables() *sql.DB {
	db, err := sql.Open("sqlite3", "server/db/main.db")
	if err != nil {
		fmt.Println(err)
		logs.FatalLog.Fatalln("Error opening database:", err)
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		fmt.Println(err)
		logs.FatalLog.Fatalln("driver instance error:", err)
	}

	absPath, err := filepath.Abs("server/sql/migrations")
	if err != nil {
		fmt.Println(err)
		logs.FatalLog.Fatalln("unable to resolve migration path:", err)
	}

	// Convert Windows path to proper file URL format
	var fileURL string
	if runtime.GOOS == "windows" {
		// Replace backslashes with forward slashes
		absPath = strings.ReplaceAll(absPath, "\\", "/")
		// For Windows, use file:// (not file:///) and handle drive letter
		fileURL = "file://" + absPath
	} else {
		// For Unix-like systems, use file:// with absolute path
		fileURL = "file://" + absPath
	}

	fmt.Printf("Migration path: %s\n", fileURL)

	m, err := migrate.NewWithDatabaseInstance(
		fileURL,
		"sqlite3",
		driver,
	)
	if err != nil {
		fmt.Println(err)
		logs.FatalLog.Fatalln("migration instance error:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println(err)
		logs.FatalLog.Fatalln("migration failed:", err)
	}

	fmt.Println("✅ migration applied!")
	logs.InfoLog.Println("✅ Database migrations applied!")
	return db
}
