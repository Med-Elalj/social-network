package modules

import (
	"fmt"
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

	var image interface{}
	if post.Image == "" {
		image = nil
	} else {
		image = post.Image
	}

	res, err := tx.Exec(`
        INSERT INTO posts (user_id, group_id, content, image_path, privacy)
        VALUES (?, ?, ?, ?, ?)`,
		uid,
		gid,
		post.Content,
		image,
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

// anas
func userdelposts(post_id int, user_id int) error {
	res, err := DB.Exec(`
		DELETE FROM posts
		WHERE id = ? AND user_id = ?`, post_id, user_id)
	if err != nil {
		return err
	}

	rowsAffect, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffect == 0 {
		return fmt.Errorf("no rows affected, post may not exist or you may not have permission")
	}
	return nil
}

// anas
func admdelposts(post_id int, user_id int, group_id int) error {
	res, err := DB.Exec(`
		DELETE FROM posts
		WHERE id = ? AND user_id = ? AND group_id = ?`, post_id, user_id, group_id)
	if err != nil {
		return err
	}

	rowsAffect, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffect == 0 {
		return fmt.Errorf("no rows affected, post may not exist or you may not have permission")
	}
	return nil
}

// anas
func updpost(newpost structs.Post) error {
	res, err := DB.Exec(`
	update post p set p.content = ?, p.image_path = ? ,p.privacy = ? from posts where p.id = ? and p.user_id = ?`, newpost.Content, newpost.ImagePath, newpost.Privacy, newpost.ID, newpost.UserId)
	if err != nil {
		return err
	}

	rowsaffect, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowsaffect == 0 {
		return fmt.Errorf("no rows affected, post may not exist or you may not have permission to update it")
	}
	return nil
}

// anas
func Insertevent(event structs.GroupEvent, uid int) (int,error) {
	tx, err := DB.Begin()
	if err != nil {
		return 0,err
	}
	_, err = tx.Exec(`INSERT INTO events (user_id,group_id,content,title,timeof) VALUES (?,?,?,?,?,?)`, uid, event.Group_id, event.Description, event.Title, event.Timeof)
	if err != nil {
		tx.Rollback()
		return 0,err
	}
	
    	lastID, err := res.LastInsertId()
    	if err != nil {
        	tx.Rollback()
        	return 0, err
    	}

    	err = tx.Commit()
    	if err != nil {
        	return 0, err
    	}
	return lastID,nil
}


func UpdatEventResp(event_id int, uid int,respond bool) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE userevent SET respond = ? WHERE event_id = ? AND user_id = ?`, respond, event_id, uid)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func InsertUserEvent(event_id int, uid int,respond bool) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO userevent (user_id, event_id, respond) VALUES (?,?,?)`, uid, event_id,respond)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func InsertGroup(gp structs.Group, uid int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	var Privacy int
	if gp.Privacy == "public" {
		Privacy = 1
	} else {
		Privacy = 0
	}
	res, err := tx.Exec(`INSERT INTO profile (display_name,avatar,description,is_public, is_user) VALUES (?,?,?,?, 0)`, gp.GroupName, gp.Avatar, gp.About, Privacy)
	if err != nil {
		tx.Rollback()
		return err
	}

	ID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`INSERT INTO "group" (id, creator_id)
	VALUES (?, ?)`,
		ID, uid)
	if err != nil {
		tx.Rollback()
		return err
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
		return err
	}

	return nil
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
