package modules

import (
	"errors"
	"fmt"

	"social-network/app/logs"
)

func InsertRequest(senderId, receiverId, typeId int) error {
	if typeId == 1 {
		fmt.Println(senderId, receiverId)
		_, err := DB.Exec(`
	  	INSERT INTO request (sender_id, receiver_id, target_id, type)
	  	SELECT ?, g.creator_id, g.id, 1
	  	FROM "group" g
	  	WHERE g.id = ?`, senderId, receiverId)
		if err != nil {
			logs.ErrorLog.Printf("error inserting new request: %q", err.Error())
			return errors.New("error inserting new request")
		}
	}
	return nil
}
