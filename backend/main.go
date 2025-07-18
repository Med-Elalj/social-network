package main

import (
	"fmt"
	"net/http"
	"time"

	routes "social-network/app"
	"social-network/app/logs"
	database "social-network/app/modules"

	_ "github.com/mattn/go-sqlite3"
)

var mux = routes.SetupMux()

func main() {
	database.DB = database.SetTables()
	defer database.DB.Close()

	fmt.Println("Server is running at http://localhost:8080")

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ErrorLog:          logs.ErrorLog,
	}

	logs.FatalLog.Fatalln(srv.ListenAndServe())
}
