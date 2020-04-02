package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type WelcomeMessage struct {
	GuildId int64  `gorm:"column:GUILDID"`
	Message string `gorm:"column:MESSAGE;type:text"`
}

func (WelcomeMessage) TableName() string {
	return "welcomemessages"
}

func UpdateWelcomeMessage(guildId int64, message string) {
	database.Database.Where(&WelcomeMessage{GuildId: guildId}).Assign(&WelcomeMessage{Message: message}).FirstOrCreate(&WelcomeMessage{})
}

func GetWelcomeMessage(guildId int64) string {
	message := WelcomeMessage{Message: "No message specified"}
	database.Database.Where(&WelcomeMessage{GuildId: guildId}).First(&message)

	return message.Message
}
