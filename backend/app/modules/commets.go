package modules

import (
	"social-network/app/structs"
	"social-network/server/logs"
)

func InsertComment(comment structs.CommentInfo, uid int) bool {
	tx, err := DB.Begin()
	if err != nil {
		logs.FatalLog.Fatalln("Database transaction error:", err)
		return false
	}

	query := `
		INSERT
			INTO comments
			(post_id, user_id, content,image_path)
		VALUES
			(?, ?, ?, ?);
	`

	res, err := tx.Exec(query,
		comment.PostID,
		uid,
		comment.Content,
		comment.ImagePath)
	if err != nil {
		logs.ErrorLog.Println("Database insertion error:", err)
		return false
	}
	lastInsertID, _ := res.LastInsertId()
	logs.InfoLog.Println("Comment inserted successfully for post:", lastInsertID)

	return true
}
