package table

import "github.com/TicketsBot/GoPanel/database"

type PermissionNode struct {
	GuildId   int64 `gorm:"column:GUILDID"`
	UserId    int64 `gorm:"column:USERID"`
	IsSupport bool  `gorm:"column:ISSUPPORT"`
	IsAdmin   bool  `gorm:"column:ISADMIN"`
}

func (PermissionNode) TableName() string {
	return "permissions"
}

func GetAdminGuilds(userId int64) []int64 {
	var nodes []PermissionNode
	database.Database.Where(&PermissionNode{UserId: userId}).Find(&nodes)

	ids := make([]int64, 0)
	for _, node := range nodes {
		ids = append(ids, node.GuildId)
	}

	return ids
}

func IsSupport(guildId int64, userId int64) bool {
	var node PermissionNode
	database.Database.Where(&PermissionNode{GuildId: guildId, UserId: userId}).Take(&node)
	return node.IsSupport
}

func IsAdmin(guildId int64, userId int64) bool {
	var node PermissionNode
	database.Database.Where(&PermissionNode{GuildId: guildId, UserId: userId}).Take(&node)
	return node.IsAdmin
}

func IsStaff(guildId int64, userId int64) bool {
	var node PermissionNode
	database.Database.Where(&PermissionNode{GuildId: guildId, UserId: userId}).Take(&node)
	return node.IsAdmin || node.IsSupport
}
