package modules

import (
	"social-network/app/logs"
	"social-network/app/structs"
)

func InsertComment(comment structs.CommentInfo, uid int) bool {
	tx, err := DB.Begin()
	if err != nil {
		logs.FatalLog.Fatalln("Database transaction error:", err)
		return false
	}

	res, err := tx.Exec(`
		INSERT INTO comments (post_id, user_id, content, image_path)
		VALUES (?, ?, ?, ?);`,
		comment.PostID,
		uid,
		comment.Content,
		comment.Image)

	if err != nil {
		tx.Rollback()
		logs.ErrorLog.Println("Database insertion error:", err)
		return false
	}

	err = tx.Commit()
	if err != nil {
		logs.ErrorLog.Println("Transaction commit error:", err)
		return false
	}

	lastInsertID, _ := res.LastInsertId()
	logs.InfoLog.Println("Comment inserted successfully for post:", lastInsertID)

	return true
}

func GetComments(commentData structs.CommentGet, uid int) ([]structs.Comments, bool) {
	tx, err := DB.Begin()
	if err != nil {
		logs.FatalLog.Fatalln("Database transaction error:", err)
		return nil, false
	}

	var comments []structs.Comments
	rows, err := tx.Query(`
        SELECT
            c.id,
            u.display_name,
            u.avatar,
            c.content,
            c.created_at,
				c.image_path,
            (
                SELECT
                    COUNT(*)
                FROM
                    likes l
                WHERE
                    l.comment_id = c.id
            ) AS like_count,
            CASE
                WHEN EXISTS (
                    SELECT
                        1
                    FROM
                        likes l
                    WHERE
                        l.user_id = ?
                        AND l.comment_id = c.id
                        AND l.post_id IS NULL
                ) THEN 1
                ELSE 0
            END AS is_liked
        FROM
            comments c
            JOIN profile u ON c.user_id = u.id
        WHERE
            c.post_id = ?
            AND (
                (
                    ? = 0
                    AND 1 = 1
                ) -- if start_id = 0, no filter on c.id
                OR (
                    c.id <= ?
                    AND ? <> 0
                ) -- if start_id != 0, filter c.id <= start_id
            )
        ORDER BY
            c.created_at DESC
        LIMIT
            10;
    `, uid, commentData.Post_id, commentData.Start, commentData.Start, commentData.Start)
	if err != nil {
		logs.ErrorLog.Printf("Error getting comments: %q", err.Error())
		return nil, false
	}
	defer rows.Close()
	for rows.Next() {
		var comment structs.Comments
		err := rows.Scan(&comment.ID, &comment.Author, &comment.AvatarUser, &comment.Content, &comment.CreatedAt, &comment.ImagePath, &comment.LikeCount, &comment.IsLiked)
		if err != nil {
			logs.ErrorLog.Printf("Error scanning comment: %q", err.Error())
			return nil, false
		}
		comments = append(comments, comment)
	}
	return comments, true
}
