package modules

import (
	"errors"
	"fmt"

	"social-network/app/logs"
	"social-network/app/structs"
)

func InsertRequest(senderId, receiverId, target, typeId int) error {
	var isToAdmin bool
	fmt.Println(target)
	if typeId == 1 && receiverId == 0 {
		isToAdmin = true
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

	notifTosent, err := UserInfoForNotification(senderId, receiverId, target)
	if err != nil {
		logs.ErrorLog.Printf("error getting user info for notification: %v", err)
		return fmt.Errorf("error getting user info for notification: %w", err)
	}

	switch typeId {
	case 1:
		if isToAdmin {
			notifTosent.Message = notifTosent.Sender.DisplayName + " wants to join " + notifTosent.Target.DisplayName + " group"
		} else {
			notifTosent.Message = notifTosent.Sender.DisplayName + " sent you a request to join " + notifTosent.Target.DisplayName + " group"
		}
		structs.NotifyUser(receiverId, "group_request", notifTosent)
	case 0:
		notifTosent.Message = notifTosent.Sender.DisplayName + " sent you a follow request"
		structs.NotifyUser(receiverId, "follow_request", notifTosent)
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
