package messagequeue

import (
	"encoding/json"
	"fmt"
)

type TicketMessage struct {
	GuildId  string `json:"guild"`
	TicketId int    `json:"ticket"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

func (c *RedisClient) ListenForMessages(message chan TicketMessage) {
	pubsub := c.Subscribe("tickets:webchat:inboundmessage")

	for {
		msg, err := pubsub.ReceiveMessage(); if err != nil {
			fmt.Println(err.Error())
			continue
		}

		var decoded TicketMessage
		if err := json.Unmarshal([]byte(msg.Payload), &decoded); err != nil {
			fmt.Println(err.Error())
			continue
		}

		message<-decoded
	}
}
