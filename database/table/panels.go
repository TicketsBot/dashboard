package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type Panel struct {
	MessageId int64 `gorm:"column:MESSAGEID"`
	ChannelId int64 `gorm:"column:CHANNELID"`
	GuildId   int64 `gorm:"column:GUILDID"` // Might be useful in the future so we store it

	Title          string `gorm:"column:TITLE;type:VARCHAR(255)"`
	Content        string `gorm:"column:CONTENT;type:TEXT"`
	Colour         int    `gorm:"column:COLOUR`
	TargetCategory int64  `gorm:"column:TARGETCATEGORY"`
	ReactionEmote  string `gorm:"column:REACTIONEMOTE;type:VARCHAR(32)"`
}

func (Panel) TableName() string {
	return "panels"
}

func AddPanel(messageId, channelId, guildId int64, title, content string, colour int, targetCategory int64, reactionEmote string) {
	database.Database.Create(&Panel{
		MessageId: messageId,
		ChannelId: channelId,
		GuildId:   guildId,

		Title:          title,
		Content:        content,
		Colour:         colour,
		TargetCategory: targetCategory,
		ReactionEmote:  reactionEmote,
	})
}

func IsPanel(messageId int64, ch chan bool) {
	var count int
	database.Database.Table(Panel{}.TableName()).Where(Panel{MessageId: messageId}).Count(&count)
	ch <- count > 0
}

func GetPanelsByGuild(guildId int64, ch chan []Panel) {
	var panels []Panel
	database.Database.Where(Panel{GuildId: guildId}).Find(&panels)
	ch <- panels
}

func DeletePanel(msgId int64) {
	database.Database.Where(Panel{MessageId: msgId}).Delete(Panel{})
}
