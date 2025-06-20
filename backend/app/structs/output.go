package structs

import (
	"encoding/json"
	"errors"
	"net/http"
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

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func JsRespond(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  code,
	})
}
