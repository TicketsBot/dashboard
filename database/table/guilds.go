package table

import (
	"encoding/base64"
	"encoding/json"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
)

type GuildCache struct {
	UserId string `gorm:"column:USERID;type:varchar(20)"` // Apparently I made this a VARCHAR in the JS version
	Guilds string `gorm:"column:guilds;type:mediumtext"`
}

func (GuildCache) TableName() string {
	return "guildscache"
}

func UpdateGuilds(userId string, guilds string) {
	var cache GuildCache
	database.Database.Where(&GuildCache{UserId: userId}).Assign(&GuildCache{Guilds: guilds}).FirstOrCreate(&cache)
}

func GetGuilds(userId string) []objects.Guild {
	var cache GuildCache
	database.Database.Where(&GuildCache{UserId: userId}).First(&cache)
	decoded, err := base64.StdEncoding.DecodeString(cache.Guilds); if err != nil {
		return make([]objects.Guild, 0)
	}

	var guilds []objects.Guild
	if err := json.Unmarshal(decoded, &guilds); err != nil {
		return make([]objects.Guild, 0)
	}

	return guilds
}
