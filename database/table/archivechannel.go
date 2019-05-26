package table

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/database"
)

type ArchiveChannel struct {
	Guild  int64 `gorm:"column:GUILDID"`
	Channel int64 `gorm:"column:CHANNELID"`
}

func (ArchiveChannel) TableName() string {
	return "archivechannel"
}

func UpdateArchiveChannel(guildId int64, channelId int64) {
	fmt.Println(channelId)
	var channel ArchiveChannel
	database.Database.Where(ArchiveChannel{Guild: guildId}).Assign(ArchiveChannel{Channel: channelId}).FirstOrCreate(&channel)
}

func GetArchiveChannel(guildId int64) int64 {
	var channel ArchiveChannel
	database.Database.Where(&ArchiveChannel{Guild: guildId}).First(&channel)

	return channel.Channel
}

