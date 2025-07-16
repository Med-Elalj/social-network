package modules

import (
	"errors"
	"fmt"

	"social-network/app/logs"
)

func InsertRequest(senderId, receiverId, target, typeId int) error {
	if typeId == 1 {
		err := DB.QueryRow(`SELECT g.creator_id FROM "group" g WHERE g.id = ?`, target).Scan(&receiverId)
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

func InsertGroupRequestFromUser(senderId, targetId, receiverId int) error {
	fmt.Println(senderId, targetId, receiverId)
	_, err := DB.Exec(`
		INSERT INTO request (sender_id, receiver_id, target_id, type)
		VALUES (?,?,?,1);`, senderId, receiverId, targetId)
	if err != nil {
		logs.ErrorLog.Printf("error inserting new request: %q", err.Error())
		return errors.New("error inserting new request")
	}
	return nil
}
