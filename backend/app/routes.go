package sn

import (
	"net/http"

	"social-network/app/handlers"
	"social-network/app/security/middleware"
	"social-network/app/ws"
)

var mux = http.NewServeMux()

func SetupMux() *http.ServeMux {
	mux.Handle("/", IndexHandler())

	// auth handlers
	mux.HandleFunc("POST /api/v1/auth", middleware.Logged_IN(handlers.Islogged))
	mux.HandleFunc("POST /api/v1/auth/login", middleware.Logged_OUT(handlers.LoginHandler))
	mux.HandleFunc("POST /api/v1/auth/register", middleware.Logged_OUT(handlers.RegisterHandler))
	mux.HandleFunc("POST /api/v1/auth/logout", handlers.LogoutHandler)

	// TODO STILL WORK IN PROGRESS
	mux.HandleFunc("POST /api/v1/ws", middleware.Logged_IN(ws.HandleConnections))

	mux.HandleFunc("GET /api/v1/get/{type}", middleware.Logged_IN(handlers.GetHandler))
	mux.HandleFunc("POST /api/v1/set/{type}", middleware.Logged_IN(handlers.SetHandler))
	return mux
}
