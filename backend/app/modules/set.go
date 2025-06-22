package modules

import (
	"database/sql"

	"social-network/app/structs"
	"social-network/server/logs"
)

// Insert new user
func InsertUser(user structs.Register) error {
	user.Password.Hash()
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	var avatar sql.NullString
	if avatar.String == "" {
		avatar = sql.NullString{String: "", Valid: false}
	} else {
		avatar = sql.NullString{String: string(avatar.String), Valid: true}
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

// Insert new post
func InsertPost(post structs.PostCreate, uid int, gid interface{}) bool {
	tx, err := DB.Begin()
	if err != nil {
		logs.Fatal(err)
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
		logs.Errorf("Database insertion error: %q", err.Error())
		return false
	}

	err = tx.Commit()
	if err != nil {
		logs.Errorf("Transaction commit error: %q", err.Error())
		return false
	}

	lastInsertID, _ := res.LastInsertId()
	logs.Println("Post inserted with ID ", lastInsertID)

	return true
}

func InsertGroup(gp structs.Group) (sql.Result, error) {
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}

	res, err := tx.Exec(`INSERT INTO profile (display_name,avatar,description, is_person) VALUES (?,?,?, 0)`, gp.UserName, gp.Avatar, gp.About)
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
		ID, gp.Cid)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	ID, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec(`INSERT INTO groupmember (group_id, person_id, active)
	VALUES (?, ?, ?)`,
		ID, gp.Cid, 3)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}

func InsertUGroup(gu structs.GroupReq, uid int) bool {
	var adminid int
	var query string
	err := DB.QueryRow(`SELECT g.creator_id FROM group g WHERE g.id = ?;`, gu.Gid).Scan(adminid)
	if err != nil {
		return false
	}
	person := 0
	active := 0
	if uid == adminid {
		active++
	}
	err = DB.QueryRow(`SELECT g.person_id FROM groupmember g WHERE g.group_id = ?;`, gu.Gid).Scan(person)
	if err == sql.ErrNoRows {
		if uid == gu.Uid {
			active++
		}
		query = `
				INSERT
				INTO groupmember
				(group_id, person_id, active)
				VALUES
				(?, ?, ?)`
	} else {
		return false
	}
	if person != 0 {
		active++
	}
	_, err = DB.Exec(query,
		gu.Gid,
		gu.Uid)
	if err != nil {
		logs.Println("Database insertion error:", err)
		return false
	}
	return true
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
