package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type ArchiveChannel struct {
	Guild   uint64 `gorm:"column:GUILDID"`
	Channel uint64 `gorm:"column:CHANNELID"`
}

func (ArchiveChannel) TableName() string {
	return "archivechannel"
}

func UpdateArchiveChannel(guildId uint64, channelId uint64) {
	var channel ArchiveChannel
	database.Database.Where(ArchiveChannel{Guild: guildId}).Assign(ArchiveChannel{Channel: channelId}).FirstOrCreate(&channel)
}

func GetArchiveChannel(guildId uint64, ch chan uint64) {
	var channel ArchiveChannel
	database.Database.Where(&ArchiveChannel{Guild: guildId}).First(&channel)
	ch <- channel.Channel
}
