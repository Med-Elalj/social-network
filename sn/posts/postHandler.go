package posts

import (
	"encoding/json"
	"errors"
	"net/http"

	"social-network/server/logs"
	"social-network/sn/db"
	"social-network/sn/handlers"
	"social-network/sn/structs"
)

type comnt struct {
	ID        int
	Content   string
	CreatedAt string
	Username  string
	// LikeCount    int
	// DislikeCount int
}

func postdata(postID string) (structs.Post, error) {
	var post structs.Post
	err := db.DB.QueryRow(`
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
	rows, err := db.DB.Query(`
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
			logs.Println("Error scanning comment:", err)
			continue
		}
		com = append(com, comment)
	}
	if err := rows.Err(); err != nil {
		logs.Println("error iterating over rows:", err)
		return nil, err
	}
	return com, nil
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
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
	isLoggedIn := handlers.Loggedin(w, r)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"post":       post,
		"comments":   comments,
		"isLoggedIn": isLoggedIn,
	})
}

func checkPostReaction(postID string, userID int, action int) error {
	if db.DB == nil {
		logs.Println("Database connection is nil!")
		return errors.New("database connection is nil")
	}
	switch action {
	case 1:
		_, err := db.DB.Exec(`
			INSERT INTO postreaction (post_id, user_id, action)
			VALUES (?, ?, ?)
			ON CONFLICT(user_id, post_id) DO UPDATE SET action = excluded.action;
`, postID, userID, true)
		return err
	case 0:
		_, err := db.DB.Exec(`
			DELETE FROM postreaction
			WHERE post_id = ? AND user_id = ?`, postID, userID)
		return err
	case -1:
		_, err := db.DB.Exec(`
	INSERT INTO postreaction (post_id, user_id, action)
	VALUES (?, ?, ?)
	ON CONFLICT(user_id, post_id) DO UPDATE SET action = excluded.action;
`, postID, userID, false)
		return err
	}
	return nil
}
