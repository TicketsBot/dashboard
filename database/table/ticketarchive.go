package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type TicketArchive struct {
	Uuid string `gorm:"column:UUID;type:varchar(36)"`
	Guild int64 `gorm:"column:GUILDID"`
	User int64 `gorm:"column:USERID"`
	Username string `gorm:"column:USERNAME;type:varchar(32)"`
	TicketId int `gorm:"column:TICKETID"`
	CdnUrl string `gorm:"column:CDNURL;type:varchar(100)"`
}

func (TicketArchive) TableName() string {
	return "ticketarchive"
}

func GetTicketArchives(guildId int64) []TicketArchive {
	var archives []TicketArchive
	database.Database.Where(&TicketArchive{Guild: guildId}).Find(&archives)

	return archives
}

func GetFilteredTicketArchives(guildId int64, userId int64, username string, ticketId int) []TicketArchive {
	var archives []TicketArchive

	query := database.Database.Where(&TicketArchive{Guild: guildId})
	if userId != 0 {
		query = query.Where(&TicketArchive{User: userId})
	}
	if username != "" {
		query = query.Where(&TicketArchive{Username: username})
	}
	if ticketId != 0 {
		query = query.Where(&TicketArchive{TicketId: ticketId})
	}

	query.Find(&archives)

	return archives
}
