package table

import (
	"github.com/TicketsBot/GoPanel/database"
	"time"
)

type ModMailArchive struct {
	Uuid      string    `gorm:"column:UUID;type:varchar(36);unique;primary_key"`
	Guild     uint64    `gorm:"column:GUILDID"`
	User      uint64    `gorm:"column:USERID"`
	CloseTime time.Time `gorm:"column:CLOSETIME"`
}

func (ModMailArchive) TableName() string {
	return "modmail_archive"
}

func (m *ModMailArchive) Store() {
	database.Database.Create(m)
}

func GetModmailArchive(uuid string, ch chan ModMailArchive) {
	var row ModMailArchive
	database.Database.Where(ModMailArchive{Uuid: uuid}).Take(&row)
	ch <- row
}

func GetModmailArchivesByUser(userId, guildId uint64, ch chan []ModMailArchive) {
	var rows []ModMailArchive
	database.Database.Where(ModMailArchive{User: userId, Guild: guildId}).Order("CLOSETIME desc").Find(&rows)
	ch <- rows
}

func GetModmailArchivesByGuild(guildId uint64, ch chan []ModMailArchive) {
	var rows []ModMailArchive
	database.Database.Where(ModMailArchive{Guild: guildId}).Order("CLOSETIME desc").Find(&rows)
	ch <- rows
}
