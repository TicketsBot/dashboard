package redis

import (
	"encoding/json"
	"github.com/apex/log"
)

type TicketCloseMessage struct {
	GuildId  uint64
	TicketId int
	User     uint64
	Reason   string
}

func (c *RedisClient) PublishTicketClose(guildId uint64, ticketId int, userId uint64, reason string) {
	settings := TicketCloseMessage{
		GuildId:  guildId,
		TicketId: ticketId,
		User:     userId,
		Reason:   reason,
	}

	encoded, err := json.Marshal(settings)
	if err != nil {
		log.Error(err.Error())
		return
	}

	c.Publish(DefaultContext(), "tickets:close", string(encoded))
}
