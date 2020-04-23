package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type PingEveryone struct {
	GuildId      uint64 `gorm:"column:GUILDID"`
	PingEveryone bool   `gorm:"column:PINGEVERYONE;type:TINYINT"`
}

func (PingEveryone) TableName() string {
	return "pingeveryone"
}

// tldr I hate gorm
func UpdatePingEveryone(guildId uint64, pingEveryone bool) {
	var settings []PingEveryone
	database.Database.Where(&PingEveryone{GuildId: guildId}).Find(&settings)

	updated := PingEveryone{guildId, pingEveryone}

	if len(settings) == 0 {
		database.Database.Create(&updated)
	} else {
		database.Database.Table("pingeveryone").Where("GUILDID = ?", guildId).Update("PINGEVERYONE", pingEveryone)
	}

	//database.Database.Where(&PingEveryone{GuildId: guildId}).Assign(&updated).FirstOrCreate(&PingEveryone{})
}

func GetPingEveryone(guildId uint64, ch chan bool) {
	pingEveryone := PingEveryone{PingEveryone: true}
	database.Database.Where(&PingEveryone{GuildId: guildId}).First(&pingEveryone)
	ch <- pingEveryone.PingEveryone
}
