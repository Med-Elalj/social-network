package db

import (
	"database/sql"
	"fmt"

	"social-network/app/structs"
)

func InsertUser(user structs.Register) error {
	user.Password.Hash()
	fmt.Println(user)

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	var avatar sql.NullString
	if user.Avatar == "" {
		avatar = sql.NullString{String: "", Valid: false}
	} else {
		avatar = sql.NullString{String: string(user.Avatar), Valid: true}
	}

	var about sql.NullString
	if user.About == "" {
		about = sql.NullString{String: "", Valid: false}
	} else {
		about = sql.NullString{String: string(user.About), Valid: true}
	}

	res, err := tx.Exec(`INSERT INTO profile (display_name,avatar,description, is_person) VALUES (?,?,?, 1)`, user.UserName, avatar, about)
	if err != nil {
		tx.Rollback()
		return err
	}

	profileID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`INSERT INTO person (id, email, first_name, last_name, password_hash, date_of_birth, gender)
	VALUES (?, ?, ?, ?, ?, ?, ?)`,
		profileID, user.Email, user.Fname, user.Lname, user.Password, user.Birthdate, user.Gender)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// func InsertPost(post structs.PostCreate, uid int) bool {
// 	tx, err := DB.Begin()
// 	if err != nil {
// 		logs.Fatal(err)
// 		return false
// 	}
// 	// Exec insert, use Exec or QueryRow with RETURNING id to get pid
// 	res, err := tx.Exec(`INSERT INTO posts (user_id, group_id, title, content, privacy) VALUES (?, ?, ?, ?, ?)`,
// 		uid,
// 		nil, // post.GroupID, // TODO: handle group posts
// 		post.Title,
// 		post.Content,
// 		0, // post.Privacy, // TODO: handle privacy
// 	)
// 	if err != nil {
// 		tx.Rollback()
// 		logs.Errorf("Database insertion error: %q", err.Error())
// 		return false
// 	}
// 	// Get the last inserted id (pid) — SQLite example
// 	pid, err := res.LastInsertId()
// 	if err != nil {
// 		tx.Rollback()
// 		logs.Errorf("Failed to get last insert id: %q", err.Error())
// 		return false
// 	}
// 	// Prepare statement for categories insert
// 	stmt, err := tx.Prepare("INSERT INTO categories (pid, category) VALUES (?, ?)")
// 	if err != nil {
// 		tx.Rollback()
// 		logs.Fatal(err)
// 		return false
// 	}
// 	defer stmt.Close()
// 	// Insert categories
// 	for _, category := range post.Categories {
// 		if _, err := stmt.Exec(pid, category); err != nil {
// 			tx.Rollback()
// 			logs.Fatal(err)
// 			return false
// 		}
// 	}
// 	if err := tx.Commit(); err != nil {
// 		logs.Fatal(err)
// 		return false
// 	}
// 	return true
// }

// func InsertComment(comment structs.CommentInfo, uid int) bool {
// 	query := `
// INSERT
// 	INTO comments
// 	(post_id, uid, content)
// VALUES
// 	(?, ?, ?)`
// 	_, err := DB.Exec(query,
// 		comment.PostID,
// 		uid,
// 		comment.Content)
// 	if err != nil {
// 		logs.Println("Database insertion error:", err)
// 		return false
// 	}
// 	return true
// }
