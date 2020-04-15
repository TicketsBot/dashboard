package table

import (
	"encoding/base64"
	"encoding/json"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/rxdn/gdl/objects/guild"
)

type GuildCache struct {
	UserId uint64 `gorm:"column:USERID"`
	Guilds string `gorm:"column:guilds;type:mediumtext"`
}

func (GuildCache) TableName() string {
	return "guildscache"
}

// this is horrible
func UpdateGuilds(userId uint64, guilds string) {
	var cache GuildCache
	database.Database.Where(&GuildCache{UserId: userId}).Assign(&GuildCache{Guilds: guilds}).FirstOrCreate(&cache)
}

func GetGuilds(userId uint64) []guild.Guild {
	var cache GuildCache
	database.Database.Where(&GuildCache{UserId: userId}).First(&cache)
	decoded, err := base64.StdEncoding.DecodeString(cache.Guilds)
	if err != nil {
		return nil
	}

	var guilds []guild.Guild
	if err := json.Unmarshal(decoded, &guilds); err != nil {
		return nil
	}

	return guilds
}
