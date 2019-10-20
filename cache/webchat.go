package cache

import (
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/manage"
)

type TicketMessage struct {
	GuildId  string `json:"guild"`
	TicketId int    `json:"ticket"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

func (c *RedisClient) ListenForMessages() {
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

		manage.SocketsLock.Lock()
		for _, socket := range manage.Sockets {
			if socket.Guild == decoded.GuildId && socket.Ticket == decoded.TicketId {
				if err := socket.Ws.WriteJSON(decoded); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		manage.SocketsLock.Unlock()
	}
}
