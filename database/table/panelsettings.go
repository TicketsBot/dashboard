package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type PanelSettings struct {
	GuildId uint64 `gorm:"column:GUILDID"`
	Title   string `gorm:"column:TITLE;type:VARCHAR(255)"`
	Content string `gorm:"column:CONTENT;type:TEXT"`
	Colour  int    `gorm:"column:COLOUR`
}

func (PanelSettings) TableName() string {
	return "panelsettings"
}

func UpdatePanelSettings(guildId uint64, title string, content string, colour int) {
	settings := PanelSettings{
		Title:   title,
		Content: content,
		Colour:  colour,
	}

	database.Database.Where(&PanelSettings{GuildId: guildId}).Assign(&settings).FirstOrCreate(&PanelSettings{})
}

func UpdatePanelTitle(guildId uint64, title string) {
	settings := PanelSettings{
		Title: title,
	}

	database.Database.Where(&PanelSettings{GuildId: guildId}).Assign(&settings).FirstOrCreate(&PanelSettings{})
}

func UpdatePanelContent(guildId uint64, content string) {
	settings := PanelSettings{
		Content: content,
	}

	database.Database.Where(&PanelSettings{GuildId: guildId}).Assign(&settings).FirstOrCreate(&PanelSettings{})
}

func UpdatePanelColour(guildId uint64, colour int) {
	settings := PanelSettings{
		Colour: colour,
	}

	database.Database.Where(&PanelSettings{GuildId: guildId}).Assign(&settings).FirstOrCreate(&PanelSettings{})
}

func GetPanelSettings(guildId uint64) PanelSettings {
	settings := PanelSettings{
		Title:   "Open A Ticket",
		Content: "React with :envelope_with_arrow: to open a ticket",
		Colour:  2335514,
	}
	database.Database.Where(PanelSettings{GuildId: guildId}).First(&settings)

	return settings
}
