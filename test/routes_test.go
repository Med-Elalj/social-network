package routes__test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"social-network/sn"
	"social-network/sn/db"
	"social-network/sn/ws"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	sqlContent, err := os.ReadFile("./fake_db.sql")
	if err != nil {
		t.Fatalf("Error reading SQL file: %v", err)
	}

	_, err = dB.Exec(string(sqlContent))
	if err != nil {
		t.Fatalf("Error executing SQL: %v", err)
	}
	t.Log("\n\n\t\tDatabase successfully created!\n\n ")
	db.SetDb(dB)
	t.Cleanup(func() {
		dB.Close()
	})
	return dB
}

// key value pairs
type kvp map[string]string

func TestUserHandler(t *testing.T) {
	setupTestDB(t)
	dummyHub := &ws.Hub{}
	mux := sn.SetupMux(dummyHub)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			for k, v := range tt.cookies {
				req.AddCookie(&http.Cookie{
					Name:  k,
					Value: v,
				})
			}

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			res := w.Result()

			if res.StatusCode != tt.eStatus {
				t.Errorf("expected status %d, got %d,%v", tt.eStatus, res.StatusCode, res)
			}

			for ghn, ghv := range res.Header {
				if ehv, exist := tt.eHeaders[ghn]; exist && ehv != ghv[0] {
					t.Errorf(`expected headers : \n\n\t\t'%v'\n\n\t\tgot '%v'`, tt.eHeaders, res.Header)
					break
				}
			}
			body := w.Body.String()
			if body != tt.eBody {
				t.Errorf("\n\n\t\texpected body '%s'\n\n\t\t\t  got '%s'", tt.eBody, body)
			}
			if tt.dbres != nil {
				data, err := dumpRows(db.DB, tt.dbres[0])
				if err != nil {
					t.Errorf("expected dbres '%v', got error '%v'", tt.dbres, err)
				}

				mismatches := compareStringSlices(tt.dbres, data)
				if len(mismatches) > 0 {
					t.Errorf("Found %d mismatched rows:", len(mismatches))
					for i, pair := range mismatches {
						t.Errorf("Mismatch %d:\n\tExpected: %q\n\tActual:   %q", i+1, pair[0], pair[1])
					}
				}
			}
		})
	}
}

// 'map[Content-Type:application/json]'
// 'map[Content-Type:[text/plain; charset=utf-8] X-Content-Type-Options:[nosniff]]'
func dumpRows(db *sql.DB, query string, args ...interface{}) ([]string, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []string
	results = append(results, query)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		var sb strings.Builder
		for i, col := range columns {
			val := values[i]
			var v interface{}
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = val
			}
			sb.WriteString(fmt.Sprintf("%s=%v", col, v))
			if i < len(columns)-1 {
				sb.WriteString(", ")
			}
		}
		results = append(results, sb.String())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func compareStringSlices(expected, actual []string) [][2]string {
	var mismatches [][2]string
	maxLen := len(expected)
	if len(actual) > maxLen {
		maxLen = len(actual)
	}

	for i := 0; i < maxLen; i++ {
		var exp, act string
		if i < len(expected) {
			exp = expected[i]
		}
		if i < len(actual) {
			act = actual[i]
		}
		if exp != act {
			mismatches = append(mismatches, [2]string{exp, act})
		}
	}

	return mismatches
}
