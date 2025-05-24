package main

import (
	"net/http"
	"slices"
)

var (
	mux   = http.NewServeMux()
	files = http.Dir("front-end/")
)

func init() {
	mux.Handle("/files/", http.StripPrefix("/files/", http.HandlerFunc(noNavigation)))

	mux.HandleFunc("/", homeHandler)

	mux.HandleFunc("/api/v1/conversations", conversationsHandler)
	mux.HandleFunc("/api/v1/mark-read", markReadHandler)
	mux.HandleFunc("/ws", hub.HandleWebSocket)

	mux.HandleFunc("/api/v1/auth", islogged)
	mux.HandleFunc("/api/v1/register", registerHandler)
	mux.HandleFunc("/api/v1/login", loginHandler)
	mux.HandleFunc("/api/v1/logout", logoutHandler)
	mux.HandleFunc("/api/v1/post", postHandler)
	mux.HandleFunc("/api/v1/category-posts", categoryPostsHandler)
	mux.HandleFunc("/api/v1/follow", forunf)
	mux.HandleFunc("/api/v1/upload", uploadHandler)
	mux.HandleFunc("GET /api/v1/get/{type}", getRouter)
	mux.HandleFunc("POST /api/v1/set/{type}", setRouter)
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
		profileHandler(w, r)
	case "posts":
		postsHandler(w, r)
	case "categories":
		categoriesHandler(w, r)
	case "conversations":
		conversationsHandler(w, r)
	case "messages":
		messagesHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func setRouter(w http.ResponseWriter, r *http.Request) {
	switch r.PathValue("type") {
	case "profile":
		profileHandler(w, r)
	case "posts":
		addPostHandler(w, r)
	case "comments":
		addCommentHandler(w, r)
	case "follow":
		forunf(w, r)
	default:
		http.NotFound(w, r)
	}
}
