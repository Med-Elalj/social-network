package comments

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"social-network/server/logs"
	"social-network/sn/db"
	"social-network/sn/handlers"
)

func checkCommentReaction(commentID string, userID int, action int) error {
	if db.DB == nil {
		logs.Println("Database connection is nil!")
		return errors.New("database connection is nil")
	}
	switch action {
	case 1:
		_, err := db.DB.Exec(`
			INSERT INTO commentreaction (comment_id, user_id, action)
			VALUES (?, ?, ?)
			ON CONFLICT(user_id, comment_id) DO UPDATE SET action = excluded.action;
`, commentID, userID, true)
		return err
	case 0:
		_, err := db.DB.Exec(`
			DELETE FROM commentreaction
			WHERE comment_id = ? AND user_id = ?`, commentID, userID)
		return err
	case -1:
		_, err := db.DB.Exec(`
	INSERT INTO commentreaction (comment_id, user_id, action)
	VALUES (?, ?, ?)
	ON CONFLICT(user_id, comment_id) DO UPDATE SET action = excluded.action;
`, commentID, userID, false)
		return err
	}
	return nil
}

type Comtresp struct {
	Success      bool   `json:"success"`
	Username     string `json:"username"`
	Content      string `json:"content"`
	CreatedAt    string `json:"createdAt"`
	CommentID    int    `json:"commentId"`
	LikeCount    int    `json:"likeCount"`
	DislikeCount int    `json:"dislikeCount"`
}

func getcomntid() (int, error) {
	var id int
	err := db.DB.QueryRow("SELECT id FROM comments ORDER BY id DESC LIMIT 1").Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("jjyaaaaaaaaaaaaaaaaaaaj")
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("userId")
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID, err := handlers.Fetchuid(cookie.Value)
	if err != nil {
		http.Error(w, "something wrong happened", http.StatusInternalServerError)
		fmt.Println("fetchuid function returned error")
		return
	}

	var comment struct {
		Postid string `json:"post_id"`
		Commet string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	// postID := r.FormValue("post_id")
	// commentContent := r.FormValue("comment")

	if comment.Postid == "" || comment.Commet == "" || userID == 0 {
		http.Error(w, "Invalid comment or user", http.StatusBadRequest)
		return
	}
	pstid, err := strconv.Atoi(comment.Postid)
	if err != nil {
		http.Error(w, "Invalid comment or user", http.StatusBadRequest)
		return
	}
	_, erro := db.DB.Exec(`
        INSERT INTO comments (content, post_id, user_id)
        VALUES (?, ?, ?)`, comment.Commet, pstid, userID)
	if erro != nil {
		http.Error(w, "Error adding comment", http.StatusInternalServerError)
		return
	}
	username, err := handlers.Getusername(userID)
	if err != nil {
		http.Error(w, "something wrong happened", http.StatusInternalServerError)
		return
	}
	comtid, err := getcomntid()
	if err != nil {
		http.Error(w, "something wrong happened", http.StatusInternalServerError)
		return
	}
	response := Comtresp{
		Success:   true,
		Username:  username,
		Content:   comment.Commet,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		CommentID: comtid,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
