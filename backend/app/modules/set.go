package modules

import (
	"errors"
	"strings"

	"social-network/app/logs"
	"social-network/app/structs"
)

func InsertPost(post structs.PostCreate, uid, gid int) bool {
	// Prepare optional fields before starting the transaction
	var groupId interface{}
	if gid == 0 {
		groupId = nil
	} else {
		groupId = gid
	}

	var image interface{}
	if post.Image == "" {
		image = nil
	} else {
		image = post.Image
	}

	// Start transaction after pre-processing
	tx, err := DB.Begin()
	if err != nil {
		logs.FatalLog.Fatalln("Database transaction error:", err)
		return false
	}

	// Insert into posts
	res, err := tx.Exec(`
		INSERT INTO posts (user_id, group_id, content, image_path, privacy)
		VALUES (?, ?, ?, ?, ?)`,
		uid, groupId, post.Content, image, post.Privacy,
	)
	if err != nil {
		tx.Rollback()
		logs.ErrorLog.Printf("Database insertion error: %q", err.Error())
		return false
	}

	// Check for insert ID error
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		logs.ErrorLog.Printf("Failed to get last insert ID: %q", err.Error())
		return false
	}

	// If privacy setting requires specific access
	if len(post.Privates) > 0 {
		for _, private := range post.Privates {
			_, err = tx.Exec(`
				INSERT INTO postrack (post_id, follower_id)
				VALUES (?, ?)`, lastInsertID, private.ID)
			if err != nil {
				tx.Rollback()
				logs.ErrorLog.Printf("Error inserting into postrack: %q", err.Error())
				return false
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		logs.ErrorLog.Printf("Transaction commit error: %q", err.Error())
		return false
	}

	return true
}

func UserFollow(uid int, tid int, followStatus string) (string, error) {
	var err error
	var newStatus string

	switch followStatus {
	case "follow", "follow back":
		{
			_, err = DB.Exec(`
				INSERT OR IGNORE INTO follow (follower_id, following_id)
				VALUES (?, ?)
				`, uid, tid)
		}

	case "unfollow":
		{
			_, err = DB.Exec(`
				DELETE FROM follow
				WHERE follower_id = ? AND following_id = ?
				`, uid, tid)
		}

	case "follow request":
		{
			_, err = DB.Exec(`
        		INSERT OR IGNORE INTO request (sender_id, receiver_id, target_id, type)
        		VALUES (?, ?, ?, 0)
    		`, uid, tid, tid)
		}
	case "cancel request":
		{
			_, err = DB.Exec(`
        		DELETE FROM request WHERE sender_id = ? AND receiver_id = ? AND type = 0`,
				uid, tid)
		}
	}
	if err != nil {
		return newStatus, errors.New("error inserting follow request " + err.Error())
	}

	newStatus, err = GetRelationship(uid, tid)
	if err != nil {
		return newStatus, errors.New("error getting relationship" + err.Error())
	}

	return newStatus, nil
}

func Insertevent(event structs.GroupEvent, uid int) (int, error) {
	// Move this to BEFORE the transaction
	members, err := GetMembers(event.Group_id)
	if err != nil {
		return 0, err
	}

	// Start transaction only after all reads are done
	tx, err := DB.Begin()
	if err != nil {
		return 0, err
	}

	// Insert into events
	res, err := tx.Exec(
		`INSERT INTO events (user_id,group_id,description,title,timeof) VALUES (?,?,?,?,?)`,
		uid, event.Group_id, event.Description, event.Title, event.Timeof,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert into userevents
	_, err = tx.Exec(
		`INSERT INTO userevents (user_id, event_id, respond) VALUES (?,?,?)`,
		uid, int(lastID), true,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert into request table for each group member
	query := `INSERT INTO request (sender_id, receiver_id, target_id, type) VALUES (?,?,?,?)`
	for _, member := range members {
		_, err = tx.Exec(query, uid, member.Uid, int(lastID), 2)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func InsertUserEvent(event_id int, uid int, respond bool) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO userevents (user_id, event_id, respond) VALUES (?,?,?)`, uid, event_id, respond)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func UpdatEventResp(event_id int, uid int, respond bool) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE userevents SET respond = ? WHERE event_id = ? AND user_id = ?`, respond, event_id, uid)
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

	if strings.Contains(gp.GroupName, " ") {
		gp.GroupName = strings.ReplaceAll(gp.GroupName, " ", "_")
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
		logs.ErrorLog.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

// anas
// func userdelposts(post_id int, user_id int) error {
// 	res, err := DB.Exec(`
// 		DELETE FROM posts
// 		WHERE id = ? AND user_id = ?`, post_id, user_id)
// 	if err != nil {
// 		return err
// 	}
// 	rowsAffect, err := res.RowsAffected()
// 	if err != nil {
// 		return err
// 	} else if rowsAffect == 0 {
// 		return fmt.Errorf("no rows affected, post may not exist or you may not have permission")
// 	}
// 	return nil
// }

// anas
// func admdelposts(post_id int, user_id int, group_id int) error {
// 	res, err := DB.Exec(`
// 		DELETE FROM posts
// 		WHERE id = ? AND user_id = ? AND group_id = ?`, post_id, user_id, group_id)
// 	if err != nil {
// 		return err
// 	}
// 	rowsAffect, err := res.RowsAffected()
// 	if err != nil {
// 		return err
// 	} else if rowsAffect == 0 {
// 		return fmt.Errorf("no rows affected, post may not exist or you may not have permission")
// 	}
// 	return nil
// }

// func IsFollowMe(tid, uid int) (bool, error) {
// 	var isFollower bool
// 	err := DB.QueryRow(`
//   SELECT EXISTS (
//     SELECT 1 FROM follow
//     WHERE follower_id = ? AND following_id = ?
//   )
// `, tid, uid).Scan(&isFollower)
// 	if err != nil {
// 		return false, err
// 	}
// 	// isFollower will be true if tid follows uid
// 	return isFollower, nil
// }

// anas
// func updpost(newpost structs.Post) error {
// 	res, err := DB.Exec(`
// 	update post p set p.content = ?, p.image_path = ? ,p.privacy = ? from posts where p.id = ? and p.user_id = ?`, newpost.Content, newpost.ImagePath, newpost.Privacy, newpost.ID, newpost.UserId)
// 	if err != nil {
// 		return err
// 	}
// 	rowsaffect, err := res.RowsAffected()
// 	if err != nil {
// 		return err
// 	} else if rowsaffect == 0 {
// 		return fmt.Errorf("no rows affected, post may not exist or you may not have permission to update it")
// 	}
// 	return nil
// }

// func inviteToGroup(invitedid, gid, sender_id int) error {
// 	_, err := DB.Exec(`
// 		INSERT INTO follow (follower_id, following_id, status)
// 		VALUES (?, ?, 0)
// 		ON CONFLICT (follower_id, following_id) DO NOTHING;`, invitedid, gid)
// 	if err != nil {
// 		logs.ErrorLog.Printf("Error inserting follow: %v", err)
// 		return fmt.Errorf("error inserting follow: %w", err)
// 	}
// 	_, err = DB.Exec(`
// 		insert into requests (sender_id, receiver_id, towhat, type)
// 		values (?, ?,? , 1)`, sender_id, invitedid, gid)
// 	if err != nil {
// 		logs.ErrorLog.Printf("Error inserting group invite request: %v", err)
// 		return fmt.Errorf("error inserting group invite request: %w", err)
// 	}
// 	logs.InfoLog.Printf("User %d invited to group %d", invitedid, gid)
// 	return nil
// }
