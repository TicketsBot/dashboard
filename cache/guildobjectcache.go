package cache

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/apex/log"
	"github.com/vmihailenco/msgpack"
	"time"
)

func (c *RedisClient) StoreGuild(guild objects.Guild) {
	packed, err := msgpack.Marshal(guild)
	if err != nil {
		log.Error(err.Error())
		return
	}

	key := fmt.Sprintf("ticketspanel:guilds:%s", string(packed))
	c.Set(key, string(packed), time.Hour*48)
}

func (c *RedisClient) GetGuildByID(guildId string, res chan *objects.Guild) {
	key := fmt.Sprintf("ticketspanel:guilds:%s", guildId)
	packed, err := c.Get(key).Result()

	if err != nil {
		res <- nil
	} else {
		var unpacked objects.Guild
		if err = msgpack.Unmarshal([]byte(packed), &unpacked); err != nil {
			log.Error(err.Error())
			res <- nil
		} else {
			res <- &unpacked
		}
	}
}

func (c *RedisClient) GuildExists(guildId string, res chan bool) {
	key := fmt.Sprintf("tickets:guilds:%s", guildId)

	intResult, err := c.Exists(key).Result()
	if err != nil {
		res <- false
		return
	}

	res <- intResult == 1
}
