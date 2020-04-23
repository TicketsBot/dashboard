package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type BlacklistNode struct {
	Assoc int    `gorm:"column:ASSOCID;type:int;primary_key;auto_increment"`
	Guild uint64 `gorm:"column:GUILDID"`
	User  uint64 `gorm:"column:USERID"`
}

func (BlacklistNode) TableName() string {
	return "blacklist"
}

func IsBlacklisted(guildId, userId uint64) bool {
	var count int
	database.Database.Table("blacklist").Where(&BlacklistNode{Guild: guildId, User: userId}).Count(&count)
	return count > 0
}

func AddBlacklist(guildId, userId uint64) {
	database.Database.Create(&BlacklistNode{Guild: guildId, User: userId})
}

func RemoveBlacklist(guildId, userId uint64) {
	database.Database.Where(BlacklistNode{Guild: guildId, User: userId}).Delete(BlacklistNode{})
}

func GetBlacklistNodes(guildId uint64) []BlacklistNode {
	var nodes []BlacklistNode
	database.Database.Where(&BlacklistNode{Guild: guildId}).Find(&nodes)
	return nodes
}
