package modules

import (
	"database/sql"
	"errors"

	"social-network/app/logs"
)

func InsertFollow(uid, gid int) error {
	res, err := DB.Exec(`INSERT INTO follow (follower_id, following_id, status)
VALUES (?, ?, -- follower_id, following_id
    (SELECT is_public FROM profile WHERE id = ?) -- following_id
);`, uid, gid, gid)
	if err != nil {
		logs.ErrorLog.Printf("Error inserting follow: %v", err)
		return errors.New("error inserting follow: databse error")
	}

	if c, err := res.RowsAffected(); c == 0 {
		logs.ErrorLog.Printf("Error inserting follow: %v", err)
		return errors.New("error inserting follow: carefull nothing changed" + err.Error())
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

type Relationship struct {
	IAmFollowing       bool `json:"i_am_following"`
	TheyAreFollowingMe bool `json:"they_are_following_me"`
	IRequested         bool `json:"i_requested"`
	TheyRequested      bool `json:"they_requested"`
}

func GetRelationship(meID, profileID int) (Relationship, error) {
	var rel Relationship

	query := `
	SELECT
		EXISTS (SELECT 1 FROM follow WHERE follower_id = ? AND following_id = ?) AS i_am_following,
		EXISTS (SELECT 1 FROM follow WHERE follower_id = ? AND following_id = ?) AS they_are_following_me,
		EXISTS (SELECT 1 FROM request WHERE sender_id = ? AND receiver_id = ? AND is_accept = 0 AND type = 1) AS i_requested,
		EXISTS (SELECT 1 FROM request WHERE sender_id = ? AND receiver_id = ? AND is_accept = 0 AND type = 1) AS they_requested
	`

	err := DB.QueryRow(query,
		meID, profileID,
		profileID, meID,
		meID, profileID,
		profileID, meID,
	).Scan(
		&rel.IAmFollowing,
		&rel.TheyAreFollowingMe,
		&rel.IRequested,
		&rel.TheyRequested,
	)

	return rel, err
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
