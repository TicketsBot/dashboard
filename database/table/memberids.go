package table

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/apex/log"
)

// Use an intermediary table to prevent a many-to-many relationship
type MemberId struct {
	MemberId int `gorm:"column:MEMBERID;primary_key;auto_increment"`
	GuildId  int64 `gorm:"column:GUILDID"`
	UserId   int64 `gorm:"column:USERID"`
}

func (MemberId) TableName() string {
	return "cache_memberids"
}

func GetMemberId(guildId, userId int64, ch chan *int) {
	var row MemberId
	database.Database.Where(&MemberId{GuildId: guildId, UserId: userId}).Take(&row)

	if row.MemberId == 0 {
		row = MemberId{
			GuildId: guildId,
			UserId:  userId,
		}
		if db := database.Database.Create(row).Scan(&row); db.Error != nil {
			log.Error(db.Error.Error())
			ch <- nil
			return
		}

		ch <- &row.MemberId
	} else {
		ch <- &row.MemberId
	}
}
