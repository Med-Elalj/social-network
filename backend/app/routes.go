package handlers

import (
	"net/http"

	auth "social-network/app/Auth"
	MW "social-network/app/Auth/middleware" // middleware
	"social-network/app/handlers"
	AH "social-network/app/handlers/Auth"   // auth handlers
	P "social-network/app/handlers/Profile" // profile handlers
	"social-network/app/ws"
)

var mux = http.NewServeMux()

func SetupMux() *http.ServeMux {
	mux.Handle("/", IndexHandler())

	// auth handlers
	mux.HandleFunc("GET /api/v1/auth/status", auth.CheckAuthHandler)
	mux.HandleFunc("POST /api/v1/auth/refresh", MW.AuthMiddleware(handlers.RefreshHandler))
	mux.HandleFunc("POST /api/v1/auth/login", AH.LoginHandler)
	mux.HandleFunc("POST /api/v1/auth/register", AH.RegisterHandler)
	mux.HandleFunc("POST /api/v1/auth/logout", auth.LogoutHandler)

	// profile handlers
	mux.HandleFunc("GET /api/v1/profile/{name}", MW.AuthMiddleware(P.ProfileHandler))
	mux.HandleFunc("POST /api/v1/settings/{type}", MW.AuthMiddleware(P.ProfileSettingsHandler))

	mux.HandleFunc("/api/v1/ws", MW.AuthMiddleware(ws.HandleConnections))

	mux.HandleFunc("/api/v1/get/{type}", MW.AuthMiddleware(handlers.GetHandler))
	mux.HandleFunc("POST /api/v1/set/{type}", MW.AuthMiddleware(handlers.SetHandler))
	return mux
}
