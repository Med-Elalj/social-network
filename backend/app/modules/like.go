package modules

import (
	"database/sql"

	"social-network/app/logs"
	"social-network/app/structs"
)

func LikeDeslike(LikeInfo structs.LikeInfo, uid int) bool {
	tx, err := DB.Begin()
	if err != nil {
		logs.FatalLog.Fatalln("Database transaction error:", err)
		return false
	}

	var PostID, CommentID sql.NullInt64

	if LikeInfo.EntityType == "post" {
		PostID = sql.NullInt64{Int64: int64(LikeInfo.EntityID), Valid: true}
		CommentID = sql.NullInt64{Valid: false}
	} else {
		PostID = sql.NullInt64{Valid: false}
		CommentID = sql.NullInt64{Int64: int64(LikeInfo.EntityID), Valid: true}
	}

	if LikeInfo.IsLiked {
		// Like exists -> delete it (toggle to dislike
		_, err = tx.Exec(`
			DELETE FROM likes 
			WHERE user_id = ?
			  AND (
			    (post_id = ? AND comment_id IS NULL)
			    OR
			    (post_id IS NULL AND comment_id = ?)
			  );
		`, uid, PostID, CommentID)
		if err != nil {
			tx.Rollback()
			logs.ErrorLog.Printf("Database delete error: %q", err.Error())
			return false
		}

		logs.InfoLog.Printf("Like removed for user %d", uid)
	} else {
		// Like does not exist -> insert it
		res, err := tx.Exec(`
			INSERT INTO likes (user_id, post_id, comment_id)
			VALUES (?, ?, ?);
		`, uid, PostID, CommentID)
		if err != nil {
			tx.Rollback()
			logs.ErrorLog.Printf("Database insert error: %q", err.Error())
			return false
		}

		lastInsertID, _ := res.LastInsertId()
		logs.InfoLog.Printf("Like inserted with ID %d", lastInsertID)
	}

	if err := tx.Commit(); err != nil {
		logs.ErrorLog.Printf("Transaction commit error: %q", err.Error())
		return false
	}

	return true
}
