package modules

import (
	"database/sql"

	"social-network/app/structs"
	"social-network/server/logs"
)

// Insert new post
func InsertPost(post structs.PostCreate, uid int, gid interface{}) bool {
	tx, err := DB.Begin()
	if err != nil {
		logs.FatalLog.Fatalln("Database transaction error:", err)
		return false
	}

	res, err := tx.Exec(`
        INSERT INTO posts (user_id, group_id, content, image_path, privacy)
        VALUES (?, ?, ?, ?, ?)`,
		uid,
		gid,
		post.Content,
		post.Image,
		post.Privacy,
	)
	if err != nil {
		tx.Rollback()
		logs.ErrorLog.Printf("Database insertion error: %q", err.Error())
		return false
	}

	err = tx.Commit()
	if err != nil {
		logs.ErrorLog.Printf("Transaction commit error: %q", err.Error())
		return false
	}

	lastInsertID, _ := res.LastInsertId()
	logs.InfoLog.Println("Post inserted with ID ", lastInsertID)

	return true
}

func InsertGroup(gp structs.Group, uid int) (sql.Result, error) {
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}

	res, err := tx.Exec(`INSERT INTO profile (display_name,avatar,description, is_user) VALUES (?,?,?, 0)`, gp.UserName, gp.Avatar, gp.About)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	ID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res, err = tx.Exec(`INSERT INTO group (id, creator_id)
	VALUES (?, ?)`,
		ID, uid)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// _, err = res.LastInsertId()
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }
	// _, err = tx.Exec(`INSERT INTO groupmember (group_id, user_id, active)
	// VALUES (?, ?, ?)`,
	// 	ID, gp.Cid, 3)
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}

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
