package table

import "github.com/TicketsBot/GoPanel/database"

type UserCanClose struct {
	Guild   int64 `gorm:"column:GUILDID;unique;primary_key"`
	CanClose *bool `gorm:"column:CANCLOSE"`
}

func (UserCanClose) TableName() string {
	return "usercanclose"
}

func IsUserCanClose(guild int64, ch chan bool) {
	var node UserCanClose
	database.Database.Where(UserCanClose{Guild: guild}).First(&node)

	if node.CanClose == nil {
		ch <- true
	} else {
		ch <- *node.CanClose
	}
}

func SetUserCanClose(guild int64, value bool) {
	database.Database.Where(&UserCanClose{Guild: guild}).Assign(&UserCanClose{CanClose: &value}).FirstOrCreate(&UserCanClose{})
}
