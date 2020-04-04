package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type ChannelCategory struct {
	GuildId  uint64 `gorm:"column:GUILDID"`
	Category uint64 `gorm:"column:CATEGORYID"`
}

func (ChannelCategory) TableName() string {
	return "channelcategory"
}

func UpdateChannelCategory(guildId uint64, categoryId uint64) {
	database.Database.Where(&ChannelCategory{GuildId: guildId}).Assign(&ChannelCategory{Category: categoryId}).FirstOrCreate(&ChannelCategory{})
}

func GetChannelCategory(guildId uint64) uint64 {
	var category ChannelCategory
	database.Database.Where(&ChannelCategory{GuildId: guildId}).First(&category)

	return category.Category
}
