package structs

import (
	"errors"
	"log"

	"github.com/mattn/go-sqlite3"
)

func sanitizeSQLiteError(err error) error {
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			// Log the full error internally
			log.Printf("SQLite error: Code=%d, ExtendedCode=%d, Message=%s",
				sqliteErr.Code, sqliteErr.ExtendedCode, sqliteErr.Error())
			switch sqliteErr.Code {
			case sqlite3.ErrConstraint:
				// For example, unique constraint violation
				return errors.New("A record with the same unique value already exists.")

			case sqlite3.ErrBusy, sqlite3.ErrLocked:
				return errors.New("The database is currently busy. Please try again shortly.")

			case sqlite3.ErrReadonly:
				return errors.New("The database is in read-only mode. Changes are not allowed.")

			case sqlite3.ErrCorrupt:
				return errors.New("The database file is corrupted. Please contact support.")

			case sqlite3.ErrNotFound:
				return errors.New("The requested record was not found.")

			default:
				return errors.New("An unexpected database error occurred. Please try again.")
			}
		}
		// For non-SQLite errors, log full detail but send generic message
		log.Printf("Non-SQLite error: %v", err)
		return errors.New("An internal error occurred. Please try again.")
	}
	return nil
}
