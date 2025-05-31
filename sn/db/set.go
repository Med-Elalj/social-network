package db

import (
	"social-network/server/logs"
	"social-network/sn/structs"
)

func InsertPost(post structs.PostCreate, uid int) bool {
	tx, err := DB.Begin()
	if err != nil {
		logs.Fatal(err)
		return false
	}

	// Exec insert, use Exec or QueryRow with RETURNING id to get pid
	res, err := tx.Exec(`INSERT INTO posts (user_id, group_id, title, content, privacy) VALUES (?, ?, ?, ?, ?)`,
		uid,
		nil, // post.GroupID, // TODO: handle group posts
		post.Title,
		post.Content,
		0, // post.Privacy, // TODO: handle privacy
	)
	if err != nil {
		tx.Rollback()
		logs.Errorf("Database insertion error: %q", err.Error())
		return false
	}

	// Get the last inserted id (pid) — SQLite example
	pid, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		logs.Errorf("Failed to get last insert id: %q", err.Error())
		return false
	}

	// Prepare statement for categories insert
	stmt, err := tx.Prepare("INSERT INTO categories (pid, category) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		logs.Fatal(err)
		return false
	}

	defer stmt.Close()

	// Insert categories
	for _, category := range post.Categories {
		if _, err := stmt.Exec(pid, category); err != nil {
			tx.Rollback()
			logs.Fatal(err)
			return false
		}
	}

	if err := tx.Commit(); err != nil {
		logs.Fatal(err)
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
		comment.Content)
	if err != nil {
		logs.Println("Database insertion error:", err)
		return false
	}
	return true
}
