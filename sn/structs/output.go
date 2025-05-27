package structs

import (
	"errors"
	"strings"

	"github.com/mattn/go-sqlite3"
)

// func sanitizeSQLiteError(err error) error {
// 	if err != nil {
// 		if sqliteErr, ok := err.(sqlite3.Error); ok {
// 			// Log the full error internally
// 			logs.Printf("SQLite error: Code=%d, ExtendedCode=%d, Message=%s",
// 				sqliteErr.Code, sqliteErr.ExtendedCode, sqliteErr.Error())
// 			switch sqliteErr.Code {
// 			case sqlite3.ErrConstraint:
// 				// For example, unique constraint violation
// 				return errors.New("A record with the same unique value already exists.")

// 			case sqlite3.ErrBusy, sqlite3.ErrLocked:
// 				return errors.New("The database is currently busy. Please try again shortly.")

// 			case sqlite3.ErrReadonly:
// 				return errors.New("The database is in read-only mode. Changes are not allowed.")

// 			case sqlite3.ErrCorrupt:
// 				return errors.New("The database file is corrupted. Please contact support.")

// 			case sqlite3.ErrNotFound:
// 				return errors.New("The requested record was not found.")

// 			default:
// 				return errors.New("An unexpected database error occurred. Please try again.")
// 			}
// 		}
// 		// For non-SQLite errors, log full detail but send generic message
// 		logs.Printf("Non-SQLite error: %v", err)
// 		return errors.New("An internal error occurred. Please try again.")
// 	}
// 	return nil
// }

func SqlConstraint(err *error) bool {
	if *err != nil {
		if sqliteErr, ok := (*err).(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			if len(sqliteErr.Error()) > 26 && sqliteErr.Error()[:26] == "UNIQUE constraint failed: " {
				e := "sorry"
				for _, v := range strings.Fields(sqliteErr.Error()[26:]) {
					v1 := strings.Split(v, ".")
					if len(v1) > 1 {
						e += " " + v1[1]
					} else {
						e += " " + v
					}
				}
				*err = errors.New(e + " already exists")
				return true
			}
		}
	}
	return false
}
