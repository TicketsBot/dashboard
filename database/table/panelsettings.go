package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type PanelSettings struct {
	GuildId int64 `gorm:"column:GUILDID"`
	Title string `gorm:"column:TITLE;type:VARCHAR(255)"`
	Content string `gorm:"column:CONTENT;type:TEXT"`
	Colour int `gorm:"column:COLOUR`
}

func (PanelSettings) TableName() string {
	return "panelsettings"
}

func UpdatePanelSettings(guildId int64, title string, content string, colour int) {
	settings := PanelSettings{
		Title: title,
		Content: content,
		Colour: colour,
	}

	database.Database.Where(&PanelSettings{GuildId: guildId}).Assign(&settings).FirstOrCreate(&PanelSettings{})
}

func GetPanelSettings(guildId int64) PanelSettings {
	settings := PanelSettings{
		Title: "Open A Ticket",
		Content: "React with :envelope_with_arrow: to open a ticket",
		Colour: 2335514,
	}
	database.Database.Where(PanelSettings{GuildId: guildId}).First(&settings)

	return settings
}
