package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var hub *Hub

func main() {
	hub = NewHub()
	go hub.Run()
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/api/conversations", conversationsHandler)
	http.HandleFunc("/api/mark-read", markReadHandler)
	http.HandleFunc("/ws", hub.HandleWebSocket)
	http.HandleFunc("/api/auth", islogged)
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/logout", logoutHandler)
	http.HandleFunc("/api/posts", postsHandler)
	http.HandleFunc("/api/add-post", addPostHandler)
	http.HandleFunc("/api/post", postHandler)
	http.HandleFunc("/api/add-comment", addCommentHandler)
	http.HandleFunc("/api/categories", categoriesHandler)
	http.HandleFunc("/api/profile", profileHandler)
	http.HandleFunc("/api/messages", messagesHandler)
	http.HandleFunc("/api/category-posts", categoryPostsHandler)
	http.HandleFunc("/api/follow", forunf)
	http.HandleFunc("/api/upload", uploadHandler)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "File upload error: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Create uploads directory with explicit permissions
    if err := os.MkdirAll("uploads", 0755); err != nil {
        http.Error(w, "Could not create directory", http.StatusInternalServerError)
        return
    }

    // Secure the filename and create path
    filename := filepath.Base(handler.Filename) // Prevent directory traversal
    filePath := filepath.Join("uploads", filename)

    dst, err := os.Create(filePath)
    if err != nil {
        http.Error(w, "Unable to save file: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        http.Error(w, "Error copying file: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Return forward-slash path for web compatibility
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "path": "/uploads/" + filename, // Note the forward slash
    })
}
func postsHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := loggedin(w, r)
	rows, err := db.Query(`
		SELECT p.id, p.title, p.content, p.created_at, u.username, 
		GROUP_CONCAT(category.name, ', ') AS categories,
		COALESCE(comment.comment_count, 0) AS comment
		FROM posts p
		INNER JOIN users u ON p.id_users = u.id
		INNER JOIN post_category ON p.id = post_category.post_id
    	INNER JOIN category ON post_category.catego_id = category.id
		    LEFT JOIN (
		select post_id, count(*) as comment_count
		from comments
		GROUP by post_id
		) as comment on comment.post_id = p.id
		GROUP BY p.id
		ORDER BY p.created_at DESC;
	`)
	if err != nil {
		http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
		log.Println("Error fetching posts:", err)
		return
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Username, &post.Categories, &post.CommentCount); err != nil {
			log.Println("Error scanning post:", err)
			continue
		}
		// fmt.Println(post.Categories)
		posts = append(posts, post)
	}
	// fmt.Println(posts)
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts":      posts,
		"isLoggedIn": isLoggedIn,
	})
}
