package sn

import (
	"net/http"
	"slices"

	"social-network/sn/comments"
	"social-network/sn/handlers"
	"social-network/sn/posts"
	"social-network/sn/ws"
)

var (
	mux   = http.NewServeMux()
	files = http.Dir("front-end/")
)

func SetupMux(hub *ws.Hub) *http.ServeMux {
	mux.Handle("/files/", http.StripPrefix("/files/", http.HandlerFunc(noNavigation)))

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/api/v1/mark-read", handlers.MarkReadHandler)
	mux.HandleFunc("/ws", hub.HandleWebSocket)

	mux.HandleFunc("/api/v1/auth", handlers.Islogged)
	mux.HandleFunc("/api/v1/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/v1/login", handlers.LoginHandler)
	mux.HandleFunc("/api/v1/logout", handlers.LogoutHandler)
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
	if slices.Contains(paths, r.URL.Path) {
		http.FileServer(files).ServeHTTP(w, r)
		return
	}
	// TODO: page
	http.NotFound(w, r)
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
