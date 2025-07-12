package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	routes "social-network/app"
	database "social-network/app/modules"
	logs "social-network/server/logs"

	_ "github.com/mattn/go-sqlite3"
)

var mux = routes.SetupMux()

func main() {
	database.DB = database.SetTables()
	// upload.EnsureUploadDir()
	defer database.DB.Close()
	fmt.Println("Server is running at https://localhost:8080")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: enableCORS(mux),
		// MaxHeaderBytes:    1024, // 1KB
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ErrorLog:          logs.ErrorLog,
	}
	logs.FatalLog.Fatalln(srv.ListenAndServe())
}

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CORS request received :", r.Method, r.URL.Path)
		// Adjust these headers as needed
		// TODO: Change to specific origins in production

		fmt.Println("request headers: ", r.Header)
		fmt.Println("origin : ", r.Header.Get("Origin"))

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("Access-Control-Allow-Headers", strings.Join(keys, ", "))
		if reqHeaders := r.Header.Get("Access-Control-Request-Headers"); reqHeaders != "" {
			w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
		} else {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		fmt.Println("Handling request for:", r.URL.Path)

		h.ServeHTTP(w, r)
	})
}
