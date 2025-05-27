package sn

import (
	"net/http"

	"social-network/sn/comments"
	"social-network/sn/handlers"
	"social-network/sn/posts"
	"social-network/sn/ws"
)

var (
	mux   = http.NewServeMux()
	files = http.Dir("front_end/")
)

func SetupMux(hub *ws.Hub) *http.ServeMux {
	mux.Handle("/front_end/", http.StripPrefix("/front_end/", http.HandlerFunc(noNavigation)))

	mux.Handle("/", http.FileServer(http.Dir("front-end/"))) // handlers.HomeHandler)
	mux.HandleFunc("/api/v1/mark-read", handlers.MarkReadHandler)
	mux.HandleFunc("/api/v1/ws", hub.HandleWebSocket)

	mux.HandleFunc("/api/v1/auth", handlers.Islogged)
	mux.HandleFunc("/api/v1/auth/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/v1/auth/login", handlers.LoginHandler)
	mux.HandleFunc("/api/v1/auth/logout", handlers.LogoutHandler)

	mux.HandleFunc("/api/v1/post", posts.PostHandler)
	mux.HandleFunc("/api/v1/category-posts", handlers.CategoryPostsHandler)
	mux.HandleFunc("/api/v1/follow", handlers.Forunf)
	mux.HandleFunc("/api/v1/upload", handlers.UploadHandler)

	mux.HandleFunc("GET /api/v1/get/{type}", getRouter)
	mux.HandleFunc("POST /api/v1/set/{type}", setRouter)
	return mux
}

var paths = []string{
	// TODO: add Front end files here
}

func noNavigation(w http.ResponseWriter, r *http.Request) {
	http.FileServer(files).ServeHTTP(w, r)
	// if slices.Contains(paths, r.URL.Path) {
	// 	return
	// }
	// // TODO: page
	// w.Write([]byte("This is a static file server. No navigation allowed.\n"))
	// http.NotFound(w, r)
}

func getRouter(w http.ResponseWriter, r *http.Request) {
	switch r.PathValue("type") {
	case "profile":
		handlers.ProfileHandler(w, r)
	case "posts":
		handlers.PostsHandler(w, r)
	case "categories":
		handlers.CategoriesHandler(w, r)
	case "conversations":
		handlers.ConversationsHandler(w, r)
	case "messages":
		handlers.MessagesHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func setRouter(w http.ResponseWriter, r *http.Request) {
	switch r.PathValue("type") {
	case "profile":
		handlers.ProfileHandler(w, r)
	case "posts":
		handlers.AddPostHandler(w, r)
	case "comments":
		comments.AddCommentHandler(w, r)
	case "follow":
		handlers.Forunf(w, r)
	default:
		http.NotFound(w, r)
	}
}
