package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type Prefix struct {
	GuildId int64  `gorm:"column:GUILDID"`
	Prefix  string `gorm:"column:PREFIX;type:varchar(8)"`
}

func (Prefix) TableName() string {
	return "prefix"
}

func UpdatePrefix(guildId int64, prefix string) {
	database.Database.Where(&Prefix{GuildId: guildId}).Assign(&Prefix{Prefix: prefix}).FirstOrCreate(&Prefix{})
}

func GetPrefix(guildId int64) string {
	prefix := Prefix{Prefix: "t!"}
	database.Database.Where(&Prefix{GuildId: guildId}).First(&prefix)

	return prefix.Prefix
}
