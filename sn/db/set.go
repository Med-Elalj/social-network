package db

import (
	"html"
	"strings"

	"social-network/server/logs"
	"social-network/sn/structs"
)

func InsertPost(post structs.PostInfo, categories []string, uid int) bool {
	query := `
	INSERT
		INTO posts
		(title, content, uid, categories)
	VALUES (?, ?, ?, ?) `
	_, err := DB.Exec(query,
		html.EscapeString(post.Title),
		html.EscapeString(post.Content),
		uid,
		strings.Join(categories, ", "))
	if err != nil {
		logs.Errorf("Database insertion error: %q", err.Error())
		return false
	}
	return true
}

func InsertComment(comment structs.CommentInfo, uid int) bool {
	query := `
INSERT
	INTO comments
	(post_id, uid, content)
VALUES
	(?, ?, ?)`
	_, err := DB.Exec(query,
		comment.PostID,
		uid,
		html.EscapeString(comment.Content))
	if err != nil {
		logs.Println("Database insertion error:", err)
		return false
	}
	return true
}
