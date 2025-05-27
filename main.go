package main

import (
	"fmt"
	"net/http"

	logs "social-network/server/logs"
	"social-network/sn"
	"social-network/sn/hub"

	_ "github.com/mattn/go-sqlite3"
)

var mux = sn.SetupMux(hub.HUB)

func main() {
	logs.InitFiles()
	SetTables()
	fmt.Println("Server is running at https://localhost:8080")
	logs.Fatal(http.ListenAndServeTLS(":8080", "server/private/cert.pem", "server/private/key.pem", mux))
}
