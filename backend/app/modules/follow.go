package modules

import (
	"database/sql"
	"errors"
	"fmt"

	"social-network/app/logs"
	"social-network/app/structs"
)

func InsertFollow(follower, following int) error {
	_, err := DB.Exec(`INSERT INTO follow (follower_id, following_id)
    VALUES (?, ?);`, follower, following)
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

	err := DB.QueryRow(`SELECT is_public FROM profile WHERE id = ?`, tid).Scan(&isPublic)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("error selecting is_public row on db: " + err.Error())
	}

	err = DB.QueryRow(`
    SELECT 
      CASE
        WHEN f1.follower_id IS NOT NULL AND f2.follower_id IS NOT NULL THEN 'follow back'
        WHEN f1.follower_id IS NOT NULL THEN 'unfollow'
        WHEN r1.sender_id IS NOT NULL THEN 'cancel request'
        WHEN r2.sender_id IS NOT NULL THEN 'accept | refuse'
        WHEN ? = 1 THEN '1'
        ELSE 'follow request'
      END AS status
    FROM 
      (SELECT follower_id FROM follow WHERE follower_id = ? AND following_id = ?) f1
    LEFT JOIN 
      (SELECT follower_id FROM follow WHERE follower_id = ? AND following_id = ?) f2
    LEFT JOIN 
      (SELECT sender_id FROM request WHERE sender_id = ? AND receiver_id = ? AND type = 0) r1
    LEFT JOIN 
      (SELECT sender_id FROM request WHERE sender_id = ? AND receiver_id = ? AND type = 0) r2
    LIMIT 1;
    `, isPublic, uid, tid, tid, uid, uid, tid, tid, uid).Scan(&status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "follow", nil
		}
		fmt.Println(err)
		return "", errors.New("error getting follow status from db: " + err.Error())
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

func DeleteRequest(senderId, receiverId, target int) error {
	_, err := DB.Exec(`
        DELETE FROM request WHERE sender_id = ? AND receiver_id = ? AND target_id = ?;`, senderId, receiverId, target)
	if err != nil {
		logs.ErrorLog.Printf("error deleting request: %q", err.Error())
		return errors.New("error deleting request")
	}
	return nil
}
