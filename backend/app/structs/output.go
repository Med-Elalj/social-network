package structs

import (
	"errors"
	"strings"

	"github.com/mattn/go-sqlite3"
)

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
