package handlers

import (
	"net/http"

	auth "social-network/app/Auth"
	MW "social-network/app/Auth/middleware" // middleware
	"social-network/app/handlers"

	// "social-network/app/handlers"
	AH "social-network/app/handlers/Auth"    // auth handlers
	PH "social-network/app/handlers/Profile" // profile handlers
	UH "social-network/app/handlers/upload"  // upload handlers
	"social-network/app/ws"
)

var mux = http.NewServeMux()

func SetupMux() http.Handler {
	mux.Handle("/", IndexHandler())

	// auth handlers
	mux.HandleFunc("GET /api/v1/auth/status", auth.CheckAuthHandler)
	mux.HandleFunc("POST /api/v1/auth/refresh", MW.AuthMiddleware(handlers.RefreshHandler))
	mux.HandleFunc("POST /api/v1/auth/login", MW.UnauthenticatedRateLimit(http.HandlerFunc(AH.LoginHandler)))
	mux.HandleFunc("POST /api/v1/auth/register", MW.UnauthenticatedRateLimit(AH.RegisterHandler))
	mux.HandleFunc("POST /api/v1/auth/logout", MW.AuthMiddleware(auth.LogoutHandler))

	// profile handlers
	mux.HandleFunc("GET /api/v1/profile/{name}", MW.AuthMiddleware((MW.AuthenticatedRateLimit(PH.ProfileHandler))))
	mux.HandleFunc("POST /api/v1/settings/{type}", MW.AuthMiddleware((MW.AuthenticatedRateLimit(PH.ProfileSettingsHandler))))
	mux.HandleFunc("GET /api/v1/followers", MW.AuthMiddleware(PH.GetFollowers))
	mux.HandleFunc("GET /api/v1/following", MW.AuthMiddleware(PH.GetFollowing))

	mux.HandleFunc("/api/v1/ws", MW.AuthMiddleware(ws.HandleConnections))

	mux.HandleFunc("/api/v1/get/{type}", MW.AuthMiddleware(MW.AuthenticatedRateLimit(handlers.GetHandler)))
	mux.HandleFunc("POST /api/v1/set/{type}", MW.AuthMiddleware(MW.AuthenticatedRateLimit(handlers.SetHandler)))
	mux.HandleFunc("POST /api/v1/upload", MW.UnauthenticatedRateLimit(UH.UploadHandler))

	fileServer := http.FileServer(http.Dir("../front-end/public/uploads"))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", fileServer))

	ServerMux := MW.CorsMiddleware(mux)
	return ServerMux
}
