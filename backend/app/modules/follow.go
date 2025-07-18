package modules

import (
	"database/sql"
	"errors"
	"fmt"

	"social-network/app/logs"
	"social-network/app/structs"
)

func InsertFollow(follower, following int) error {
	// Check if relationship already exists
	var exists bool
	err := DB.QueryRow(`SELECT EXISTS (SELECT 1 FROM follow WHERE follower_id = ? AND following_id = ?)`,
		follower, following).Scan(&exists)
	if err != nil {
		logs.ErrorLog.Printf("Error checking existing follow: %v", err)
		return errors.New("error checking existing follow relationship")
	}
	fmt.Println("is eixists", exists)
	if exists {
		// Relationship already exists, don't insert duplicate
		return nil
	}

	_, err = DB.Exec(`INSERT INTO follow (follower_id, following_id) VALUES (?, ?)`,
		follower, following)
	if err != nil {
		logs.ErrorLog.Printf("Error inserting follow: %v", err)
		return errors.New("error inserting follow: database error")
	}
	return nil
}

// DeleteFollow removes a follow relationship between a user and a group.
func DeleteFollow(uid, gid int) error {
	_, err := DB.Exec(`DELETE FROM follow WHERE follower_id = ? AND following_id = ? and status <> 2;`, uid, gid)
	if err != nil {
		logs.ErrorLog.Printf("Error deleting follow: %v", err)
		return errors.New("error deleting follow: database error")
	}
	if err == sql.ErrNoRows {
		logs.ErrorLog.Printf("No follow relationship found for user %d and group %d", uid, gid)
		return errors.New("no follow relationship found")
	}
	return err
}

// AcceptFollow updates the status of a follow relationship to accepted (1).
func AcceptFollow(uid, gid, followerID int) error {
	res, err := DB.Exec(`UPDATE follow
	SET status = 1
	WHERE follower_id = ? -- user id to accept
	AND (
	    (
	        following_id = ? -- group id
	        AND EXISTS (
	            SELECT 1 FROM 'group'
	            WHERE id = ? -- group id
	            AND creator_id = ? -- admin uid
	        )
	    )
	    OR (
	        following_id = ? -- admin user id
	        AND EXISTS (
	            SELECT 1 FROM user
	            WHERE id = ? -- admin user id
	        )
	    )
	);`, followerID, gid, gid, uid, uid, uid)
	if err != nil {
		logs.ErrorLog.Printf("Error accepting follow: %v", err)
		return errors.New("error accepting follow: database error")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		logs.ErrorLog.Printf("No follow relationship found for user %d and group %d", followerID, gid)
		return errors.New("no follow relationship found or already accepted")
	}

	return err
}

func GetUserIDByUsername(username string) (int, error) {
	var userID int
	err := DB.QueryRow("SELECT id FROM profile WHERE display_name = ?", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func GetRelationship(uid, tid int) (string, error) {
	var status string
	var isPublic bool

	// Get target user's privacy status
	err := DB.QueryRow(`SELECT is_public FROM profile WHERE id = ?`, tid).Scan(&isPublic)
	if err != nil {
		return "", fmt.Errorf("error getting user privacy: %v", err)
	}

	query := `
		SELECT CASE
			-- 1. They sent you a request (you can accept/refuse)
			WHEN EXISTS (
				SELECT 1 FROM request WHERE sender_id = ? AND receiver_id = ? AND type = 0
			) THEN 'accept | refuse'

			-- 2. They follow you and you don't follow them and they are public
			WHEN EXISTS (
				SELECT 1 FROM follow WHERE follower_id = ? AND following_id = ?
			)
			AND NOT EXISTS (
				SELECT 1 FROM follow WHERE follower_id = ? AND following_id = ?
			)
			AND ? = 1 THEN 'follow back'

			-- 3. You sent a follow request and they are private
			WHEN EXISTS (
				SELECT 1 FROM request WHERE sender_id = ? AND receiver_id = ? AND type = 0
			)
			AND ? = 0 THEN 'cancel request'

			-- 4. You follow them
			WHEN EXISTS (
				SELECT 1 FROM follow WHERE follower_id = ? AND following_id = ?
			) THEN 'unfollow'

			-- 5. No relation and they are private
			WHEN ? = 0 THEN 'follow request'

			-- 6. Default: public, no relation
			ELSE 'follow'
		END AS status
	`

	err = DB.QueryRow(query,
		tid, uid, // 1. They sent request

		tid, uid, // 2. They follow you
		uid, tid, //    You don’t follow them
		isPublic, //    They are public

		uid, tid, // 3. You sent follow request
		isPublic, //    They are private

		uid, tid, // 4. You follow them

		isPublic, // 5. They are private
	).Scan(&status)

	if err == sql.ErrNoRows {
		status = "follow"
	} else if err != nil {
		return "", fmt.Errorf("error querying relationship: %v", err)
	}

	return status, nil
}

func CreateFollow(fromID, toID int) (bool, error) {
	var isPublic bool
	err := DB.QueryRow("SELECT is_public FROM profile WHERE id = ?", toID).Scan(&isPublic)
	if err != nil {
		logs.ErrorLog.Println(err)
		return false, err
	}

	return true, nil
}

func SetUnfollow(fromID, toID int) error {
	// 1. Try to delete from `follow` first
	res, err := DB.Exec(`
		DELETE FROM follow
		WHERE follower_id = ? AND following_id = ?
	`, fromID, toID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// 2. If no follow was found, try deleting request instead
	if rowsAffected == 0 {
		_, err := DB.Exec(`
			DELETE FROM request
			WHERE sender_id = ? AND receiver_id = ? AND type = 1
		`, fromID, toID)
		if err != nil {
			logs.ErrorLog.Printf("Error deleting follow request: %v", err)
			return err
		}
	}

	return nil
}

func GetFollowRequests(uid int) ([]structs.Gusers, error) {
	var users []structs.Gusers
	rows, err := DB.Query(`
	SELECT u.id, u.display_name, u.avatar
	from users profile
	LEFT JOIN request ON receiver_id = ?
	WHERE TYPE = 'follow_request' AND is_accept = 0
	order by created_at desc;
	`, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return users, nil
		}
		logs.ErrorLog.Printf("Error getting follow requests: %v", err)
		return nil, errors.New("error getting follow requests")
	}
	defer rows.Close()
	for rows.Next() {
		var user structs.Gusers
		err := rows.Scan(&user.Uid, &user.Name, &user.Avatar)
		if err != nil {
			logs.ErrorLog.Printf("Error scanning follow request: %v", err)
			return nil, errors.New("error scanning follow request")
		}
		users = append(users, user)
	}
	return users, nil
}

func DeleteRequest(senderId, uid, target, Type int) error {
	var result sql.Result
	var err error

	if Type == 1 {
		// Delete all requests for type 1 (ignore sender)
		result, err = DB.Exec(`
            DELETE FROM request WHERE receiver_id = ? AND target_id = ? AND type = ?`,
			uid, target, Type)
	} else {
		// Delete specific request for type 0 (include sender)
		result, err = DB.Exec(`
            DELETE FROM request WHERE sender_id = ? AND receiver_id = ? AND target_id = ? AND type = ?`,
			senderId, uid, uid, Type)
	}

	if err != nil {
		logs.ErrorLog.Printf("error deleting request: %q", err.Error())
		return errors.New("error deleting request")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logs.ErrorLog.Printf("error getting rows affected: %q", err.Error())
		return errors.New("error getting rows affected")
	}
	if rowsAffected == 0 {
		logs.ErrorLog.Println("no rows affected, request not found")
		return errors.New("request not found")
	}
	return nil
}
