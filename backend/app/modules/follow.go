package modules

import (
	"database/sql"
	"errors"

	"social-network/server/logs"
)

func InsertFollow(uid, gid int) error {
	res, err := DB.Exec(`INSERT INTO follow (follower_id, following_id, status)
VALUES (?, ?, -- follower_id, following_id
    (SELECT is_public FROM profile WHERE id = ?) -- following_id
);`, uid, gid, gid)
	if err != nil {
		logs.Errorf("Error inserting follow: %v", err)
		return errors.New("error inserting follow: databse error")
	}

	if c, err := res.RowsAffected(); c == 0 {
		logs.Errorf("Error inserting follow: %v", err)
		return errors.New("error inserting follow: carefull nothing changed" + err.Error())
	}

	return nil
}

// DeleteFollow removes a follow relationship between a user and a group.
func DeleteFollow(uid, gid int) error {
	_, err := DB.Exec(`DELETE FROM follow WHERE follower_id = ? AND following_id = ? and status <> 2;`, uid, gid)
	if err != nil {
		logs.Errorf("Error deleting follow: %v", err)
		return errors.New("error deleting follow: database error")
	}
	if err == sql.ErrNoRows {
		logs.Errorf("No follow relationship found for user %d and group %d", uid, gid)
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
		logs.Errorf("Error accepting follow: %v", err)
		return errors.New("error accepting follow: database error")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		logs.Errorf("No follow relationship found for user %d and group %d", followerID, gid)
		return errors.New("no follow relationship found or already accepted")
	}

	return err
}
