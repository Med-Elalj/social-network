package main

import (
	"fmt"
	"net/http"

	logs "social-network/server/logs"
	"social-network/sn"
	"social-network/sn/db"
	"social-network/sn/upload"

	_ "github.com/mattn/go-sqlite3"
)

var mux = sn.SetupMux()

func main() {
	logs.InitFiles()
	db.SetTables()
	upload.EnsureUploadDir()
	defer db.DB.Close()
	fmt.Println("Server is running at https://localhost:8080")
	logs.Fatal(http.ListenAndServeTLS(":8080", "server/private/cert.pem", "server/private/key.pem", mux))
}
