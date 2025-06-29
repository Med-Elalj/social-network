package handlers

import (
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/handlers"
	"social-network/app/ws"
)

var mux = http.NewServeMux()

func SetupMux() *http.ServeMux {
	mux.Handle("/", IndexHandler())

	// auth handlers
	mux.HandleFunc("GET /api/v1/auth/status", auth.CheckAuthHandler)
	mux.HandleFunc("GET /api/v1/auth/refresh", handlers.RefreshHandler)
	mux.HandleFunc("POST /api/v1/auth/login", auth.LoginHandler)
	mux.HandleFunc("POST /api/v1/auth/register", auth.RegisterHandler)
	mux.HandleFunc("POST /api/v1/auth/logout", auth.LogoutHandler)

	mux.HandleFunc("/api/v1/ws", ws.HandleConnections)

	mux.HandleFunc("POST /api/v1/get/{type}", handlers.GetHandler)
	mux.HandleFunc("POST /api/v1/set/{type}", handlers.SetHandler)
	return mux
}
