package sn

import (
	"net/http"

	"social-network/sn/auth/middleware"
	"social-network/sn/handlers"
	"social-network/sn/ws"
)

var mux = http.NewServeMux()

func SetupMux() *http.ServeMux {
	// TODO fix this
	mux.Handle("/", IndexHandler())

	// TODO IMPLEMENT
	// mux.HandleFunc("/api/v1/mark-read", handlers.MarkReadHandler)

	// TODO IMPLEMENT
	mux.HandleFunc("POST /api/v1/auth", middleware.Logged_IN(handlers.Islogged))
	mux.HandleFunc("POST /api/v1/auth/register", middleware.Logged_OUT(handlers.RegisterHandler))
	mux.HandleFunc("POST /api/v1/auth/login", middleware.Logged_OUT(handlers.LoginHandler))
	mux.HandleFunc("POST /api/v1/auth/logout", handlers.LogoutHandler)
	// STILL WORK IN PROGRESS
	mux.HandleFunc("POST /api/v1/ws", middleware.Logged_IN(ws.HandleConnections))

	mux.HandleFunc("GET /api/v1/get/{type}", middleware.Logged_IN(GetHandler))
	mux.HandleFunc("POST /api/v1/set/{type}", middleware.Logged_IN(PostHandler))
	return mux
}
