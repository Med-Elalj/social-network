package modules

import (
	"errors"

	"social-network/app/logs"
)

func InsertRequest(senderId, receiverId, target, typeId int) error {

	if typeId == 1 {
		err := DB.QueryRow(`select g.creator_id from groups g where g.id = ?`, target).Scan(&receiverId)
		if err != nil {
			logs.ErrorLog.Printf("error getting group creator id: %q", err.Error())
			return errors.New("error getting group creator id")
		}
		if receiverId == 0 {
			logs.ErrorLog.Println("Group creator not found")
			return errors.New("group creator not found")
		}
	}
	_, err := DB.Exec(`
          INSERT INTO request (sender_id, receiver_id, target_id, type)
    VALUES (?, ?, ?, ?);`, senderId, receiverId, target, typeId)
	if err != nil {
		logs.ErrorLog.Printf("error inserting new request: %q", err.Error())
		return errors.New("error inserting new request")
	}
	return nil
}
