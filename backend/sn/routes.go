package sn

import (
	"net/http"

	"social-network/sn/handlers"
	"social-network/sn/security/middleware"
	"social-network/sn/ws"
)

var mux = http.NewServeMux()

func SetupMux() *http.ServeMux {
	mux.Handle("/", IndexHandler())

	mux.HandleFunc("POST /api/v1/auth", middleware.Logged_IN(handlers.Islogged))
	mux.HandleFunc("POST /api/v1/auth/register", middleware.Logged_OUT(handlers.RegisterHandler))
	mux.HandleFunc("POST /api/v1/auth/login", middleware.Logged_OUT(handlers.LoginHandler))
	mux.HandleFunc("POST /api/v1/auth/logout", handlers.LogoutHandler)
	// mux.HandleFunc("POST /api/v1/profile",middleware.Logged_IN(handlers.ProfileHandler))

	// TODO STILL WORK IN PROGRESS
	mux.HandleFunc("POST /api/v1/ws", middleware.Logged_IN(ws.HandleConnections))

	mux.HandleFunc("GET /api/v1/get/{type}", middleware.Logged_IN(GetHandler))
	mux.HandleFunc("POST /api/v1/set/{type}", middleware.Logged_IN(PostHandler))
	return mux
}
