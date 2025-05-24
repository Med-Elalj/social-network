package main

import (
	"fmt"
	"log"
	"net/http"

	"social-network/sn"
	"social-network/sn/hub"

	_ "github.com/mattn/go-sqlite3"
)

var mux = sn.SetupMux(hub.HUB)

func main() {
	fmt.Println("Server is running at https://localhost:8080")
	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", mux))
}
