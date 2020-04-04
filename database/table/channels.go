package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type Channel struct {
	ChannelId uint64 `gorm:"column:CHANNELID"`
	GuildId   uint64 `gorm:"column:GUILDID"`
	Name      string `gorm:"column:NAME;type:VARCHAR(32)"`
	Type      int    `gorm:"column:CHANNELTYPE;type:TINYINT(1)"`
}

func (Channel) TableName() string {
	return "Channel"
}

func StoreChannel(channelId, guildId uint64, name string, channelType int) {
	channel := Channel{
		ChannelId: channelId,
		GuildId:   guildId,
		Name:      name,
		Type:      channelType,
	}

	database.Database.Where(&Channel{ChannelId: channelId}).Assign(&channel).FirstOrCreate(&Channel{})
}

func DeleteChannel(channelId uint64) {
	var node Channel
	database.Database.Where(Channel{ChannelId: channelId}).Take(&node)
	database.Database.Delete(&node)
}

func GetCachedChannelsByGuild(guildId uint64, res chan []Channel) {
	var nodes []Channel
	database.Database.Where(Channel{GuildId: guildId}).Find(&nodes)
	res <- nodes
}

// Util function ig
func GetCategories(guildId uint64, res chan []Channel) {
	channelsChan := make(chan []Channel)
	go GetCachedChannelsByGuild(guildId, channelsChan)
	channels := <-channelsChan

	var categories []Channel

	for _, channel := range channels {
		if channel.Type == 4 {
			categories = append(categories, channel)
		}
	}

	res <- categories
}
