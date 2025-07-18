package modules

import (
	"fmt"

	"social-network/app/logs"
)

func AddDm(sUname, rUname int, msg string) error {
	result, err := DB.Exec(`
	INSERT INTO
	message (sender_id, receiver_id, content)
	VALUES (?, ?, ?);
	`, sUname, rUname, msg)
	if err != nil {
		logs.ErrorLog.Println("Database insertion error:", err)
		return fmt.Errorf("could not Save in database")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		logs.ErrorLog.Printf("Failed to insert message: %v, Rows affected: %d", err, rowsAffected)
		return fmt.Errorf("could not Save in database")
	}
	return nil
}
