package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type ChannelCategory struct {
	GuildId int64 `gorm:"column:GUILDID"`
	Category int64 `gorm:"column:CATEGORYID"`
}

func (ChannelCategory) TableName() string {
	return "channelcategory"
}

func UpdateChannelCategory(guildId int64, categoryId int64) {
	database.Database.Where(&ChannelCategory{GuildId: guildId}).Assign(&ChannelCategory{Category: categoryId}).FirstOrCreate(&ChannelCategory{})
}

func GetChannelCategory(guildId int64) int64 {
	var category ChannelCategory
	database.Database.Where(&ChannelCategory{GuildId: guildId}).First(&category)

	return category.Category
}
