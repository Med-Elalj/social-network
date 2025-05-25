package routes__test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"social-network/sn"
	"social-network/sn/db"
	"social-network/sn/ws"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
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
	fmt.Println("\n\n\t\tDatabase successfully created!\n\n ")
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
			req := httptest.NewRequest(tt.method, tt.path, nil)
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
				t.Errorf("expected status %d, got %d", tt.eStatus, res.StatusCode)
			}

			if res.StatusCode == tt.eStatus {
				for ghn, ghv := range res.Header {

					if ehv, exist := tt.eHeaders[ghn]; exist && ehv != ghv[0] {
						t.Errorf(`expected headers : \n\n\t\t'%v'\n\n\t\tgot '%v'`, tt.eHeaders, res.Header)
						break
					}
				}
				body := w.Body.String()
				if body != tt.eBody {
					t.Errorf("\n\n\t\texpected error body '%s'\n\n\t\t got '%s' \n%v\n%v", tt.eBody, body, []rune(tt.eBody), []rune(body))
				}
				// var got map[string]interface{}
				// json.NewDecoder(res.Body).Decode(&got)

				// expected := make(map[string]interface{})
				// json.NewDecoder(bytes.NewBufferString(tt.eBody)).Decode(&expected)

				// if got["id"] != expected["id"] || got["name"] != expected["name"] {
				// 	t.Errorf("expected body %+v, got %+v", expected, got)
				// }
			} else {
				body := w.Body.String()
				if body != tt.eBody {
					t.Errorf("expected error body %q, got %q", tt.eBody, body)
				}
			}
		})
	}
}

var tests = []struct {
	name     string
	method   string
	path     string
	headers  kvp
	cookies  kvp
	eStatus  int
	eHeaders kvp
	eBody    string
}{
	{"GET root 1   ", "GET", "/", kvp{"Authorization": "Bearer testtoken"}, kvp{},
		http.StatusOK, kvp{}, `hello world`},

	{"Isloged in false no uId  ",
		"POST", "/api/v1/auth", kvp{"Authorization": "Bearer testtoken"}, kvp{},
		http.StatusOK, kvp{"Content-Type": "application/json"}, `{"isLoggedIn":false}` + "\n"},

	{"Isloged in false bad cookie name ",
		"POST", "/api/v1/auth", kvp{"Authorization": "Bearer testtoken"}, kvp{"userid": "6352337196a2449cb772b524818bea36"},
		http.StatusOK, kvp{"Content-Type": "application/json"}, `{"isLoggedIn":false}` + "\n"},

	{"Isloged in false no uid value ",
		"POST", "/api/v1/auth", kvp{"Authorization": "Bearer testtoken"}, kvp{"userId": ""},
		http.StatusOK, kvp{"Content-Type": "application/json"}, `{"isLoggedIn":false}` + "\n"},

	{"Isloged in false bad uid value ",
		"POST", "/api/v1/auth", kvp{"Authorization": "Bearer testtoken"}, kvp{"userId": "6352337196a2449cb772b524818bea37"},
		http.StatusOK, kvp{"Content-Type": "application/json"}, `{"isLoggedIn":false}` + "\n"},

	{"Isloged in true   ",
		"POST", "/api/v1/auth", kvp{"Authorization": "Bearer testtoken"}, kvp{"userId": "6352337196a2449cb772b524818bea36"},
		http.StatusOK, kvp{"Content-Type": "application/json"}, `{"isLoggedIn":true}` + "\n"},

	// {"TEST NAME   ",
	// 	"METHOD", "/URL/PATH&querry=value", []kvp{{"Header Name": "Header Value"}}, kvp{"Cookie Name": "Cookie Value"},
	// 	http.ExpectedStatus, []kvp{{"Expected Header Name": "Expected Header Value"}}, `Expected Body`},

}
