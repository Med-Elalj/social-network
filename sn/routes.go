package sn

import (
	"net/http"

	"social-network/sn/auth/middleware"
	"social-network/sn/handlers"
)

var mux = http.NewServeMux()

func SetupMux() *http.ServeMux {
	// TODO fix this
	mux.HandleFunc("/", IndexHandler)

	// TODO IMPLEMENT
	// mux.HandleFunc("/api/v1/mark-read", handlers.MarkReadHandler)

	// TODO IMPLEMENT
	// mux.HandleFunc("/api/v1/auth", handlers.Islogged)
	mux.HandleFunc("/api/v1/auth/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/v1/auth/login", handlers.LoginHandler)

	// TODO add them to get and set routers
	// mux.HandleFunc("/api/v1/follow", handlers.Forunf)
	// mux.HandleFunc("/api/v1/upload", handlers.UploadHandler)

	mux.HandleFunc("GET /api/v1/get/{type}", middleware.Mdlw_router(GetHandler_in, GetHandler_out))
	mux.HandleFunc("POST /api/v1/set/{type}", middleware.MdlwLogged_IN(PostHandler))
	return mux
}
