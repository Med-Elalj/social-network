package handlers

import (
	"net/http"

	auth "social-network/app/Auth"
	"social-network/app/Auth/middleware"
	"social-network/app/handlers"
	"social-network/app/ws"
)

var mux = http.NewServeMux()

func SetupMux() *http.ServeMux {
	mux.Handle("/", IndexHandler())

	// auth handlers
	mux.HandleFunc("POST /api/v1/auth", middleware.Logged_IN(auth.Islogged))
	mux.HandleFunc("POST /api/v1/auth/login", auth.LoginHandler)
	mux.HandleFunc("POST /api/v1/auth/register", auth.RegisterHandler)
	mux.HandleFunc("POST /api/v1/auth/logout", auth.LogoutHandler)

	mux.HandleFunc("/api/v1/ws", middleware.Logged_IN(ws.HandleConnections))

	mux.HandleFunc("POST /api/v1/get/{type}", middleware.Logged_IN(handlers.GetHandler))
	mux.HandleFunc("POST /api/v1/set/{type}", middleware.Logged_IN(handlers.SetHandler))
	return mux
}
