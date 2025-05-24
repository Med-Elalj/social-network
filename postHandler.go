package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Post struct {
	ID           int
	Title        string
	Content      string
	CreatedAt    string
	Username     string
	Categories   string
	CommentCount int
	Attachement	 string
	Status	 	 string
	// LikeCount    int
	// DislikeCount int
}

type User struct {
	ID       int
	Username string
	Email    string
	fname    string
	lname    string
	Status	 string
	followers []int
	followed  []int
}

type comnt struct {
	ID        int
	Content   string
	CreatedAt string
	Username  string
	// LikeCount    int
	// DislikeCount int
}


func postdata(postID string) (Post, error) {
    var post Post
    err := db.QueryRow(`
        SELECT 
        posts.id, 
        posts.title, 
        posts.content, 
        posts.created_at,
        posts.attachement,
        posts.status,
        GROUP_CONCAT(category.name, ', ') AS categories,
        users.username AS Username
        FROM posts
        LEFT JOIN post_category ON posts.id = post_category.post_id
        LEFT JOIN category ON post_category.catego_id = category.id
        INNER JOIN users ON posts.id_users = users.id
        WHERE posts.id = ?`, postID).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Attachement, &post.Status, &post.Categories, &post.Username)

    return post, err
}


func commentdata(postID string) ([]comnt, error) {
	var com []comnt
	rows, err := db.Query(`
	SELECT comments.id, comments.content, comments.created_at, users.username
	FROM comments 
	JOIN users ON comments.user_id = users.id
	WHERE comments.post_id = ?
	ORDER BY comments.created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment comnt
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.Username); err != nil {
			log.Println("Error scanning comment:", err)
			continue
		}
		com = append(com, comment)
	}
	if err := rows.Err(); err != nil {
		log.Println("error iterating over rows:", err)
		return nil, err
	}
	return com, nil
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, `{"error": "Post ID is required"}`, http.StatusBadRequest)
		return
	}

	post, err := postdata(postID)
	if err != nil {

		http.Error(w, `{"error": "Post not found"}`, http.StatusNotFound)
		return
	}

	comments, err := commentdata(postID)
	if err != nil {
		// fmt.Print("loooooooooool",err)
		http.Error(w, `{"error": "Error fetching comments"}`, http.StatusInternalServerError)
		return
	}
	// fmt.Println(post, comments)
	isLoggedIn := loggedin(w, r)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"post":       post,
		"comments":   comments,
		"isLoggedIn": isLoggedIn,
	})
}
