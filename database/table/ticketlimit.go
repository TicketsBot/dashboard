package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type TicketLimit struct {
	GuildId uint64 `gorm:"column:GUILDID"`
	Limit   int    `gorm:"column:TICKETLIMIT"`
}

func (TicketLimit) TableName() string {
	return "ticketlimit"
}

func UpdateTicketLimit(guildId uint64, limit int) {
	database.Database.Where(&TicketLimit{GuildId: guildId}).Assign(&TicketLimit{Limit: limit}).FirstOrCreate(&TicketLimit{})
}

func GetTicketLimit(guildId uint64) int {
	limit := TicketLimit{Limit: 5}
	database.Database.Where(&TicketLimit{GuildId: guildId}).First(&limit)

	return limit.Limit
}
