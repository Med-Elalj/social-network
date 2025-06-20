package main

import (
	"fmt"
	"net/http"
	"time"

	sn "social-network/app"
	db "social-network/app/modules"
	logs "social-network/server/logs"

	_ "github.com/mattn/go-sqlite3"
)

var mux = sn.SetupMux()

func main() {
	logs.InitFiles()
	db.SetTables()
	// upload.EnsureUploadDir()
	defer db.DB.Close()
	fmt.Println("Server is running at https://localhost:8080")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		// MaxHeaderBytes:    1024, // 1KB
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ErrorLog:          logs.GetLogger("server"),
	}
	logs.Fatal(srv.ListenAndServeTLS("./server/private/cert.pem", "./server/private/key.pem"))
}
