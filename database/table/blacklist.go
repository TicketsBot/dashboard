package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type BlacklistNode struct {
	Assoc int `gorm:"column:ASSOCID;type:int;primary_key;auto_increment"`
	Guild int64  `gorm:"column:GUILDID"`
	User int64 `gorm:"column:USERID"`
}

func (BlacklistNode) TableName() string {
	return "blacklist"
}

func IsBlacklisted(guildId, userId int64) bool {
	var count int
	database.Database.Table("blacklist").Where(&BlacklistNode{Guild: guildId, User: userId}).Count(&count)
	return count > 0
}

func AddBlacklist(guildId, userId int64) {
	database.Database.Create(&BlacklistNode{Guild: guildId, User: userId})
}

func RemoveBlacklist(guildId, userId int64) {
	database.Database.Delete(&BlacklistNode{Guild: guildId, User: userId})
}

func GetBlacklistNodes(guildId int64) []BlacklistNode {
	var nodes []BlacklistNode
	database.Database.Where(&BlacklistNode{Guild: guildId}).Find(&nodes)
	return nodes
}
