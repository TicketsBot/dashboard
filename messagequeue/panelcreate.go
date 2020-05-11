package messagequeue

import (
	"encoding/json"
	"github.com/TicketsBot/database"
	"github.com/apex/log"
)

func (c *RedisClient) PublishPanelCreate(settings database.Panel) {
	encoded, err := json.Marshal(settings); if err != nil {
		log.Error(err.Error())
		return
	}

	c.Publish("tickets:panel:create", string(encoded))
}

