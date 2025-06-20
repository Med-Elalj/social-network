package modules

import (
	"time"

	"social-network/app/structs"
)

func filter(messages []structs.Message, t time.Time) []structs.Message {
	var filteredMessages []structs.Message
	for _, message := range messages {
		if message.Time.Before(t) {
			filteredMessages = append(filteredMessages, message)
		}
	}
	return filteredMessages
}
