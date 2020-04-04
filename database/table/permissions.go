package table

import "github.com/TicketsBot/GoPanel/database"

type PermissionNode struct {
	GuildId   uint64 `gorm:"column:GUILDID"`
	UserId    uint64 `gorm:"column:USERID"`
	IsSupport bool   `gorm:"column:ISSUPPORT"`
	IsAdmin   bool   `gorm:"column:ISADMIN"`
}

func (PermissionNode) TableName() string {
	return "permissions"
}

func GetAdminGuilds(userId uint64) []uint64 {
	var nodes []PermissionNode
	database.Database.Where(&PermissionNode{UserId: userId}).Find(&nodes)

	ids := make([]uint64, 0)
	for _, node := range nodes {
		ids = append(ids, node.GuildId)
	}

	return ids
}

func IsSupport(guildId uint64, userId uint64) bool {
	var node PermissionNode
	database.Database.Where(&PermissionNode{GuildId: guildId, UserId: userId}).Take(&node)
	return node.IsSupport
}

func IsAdmin(guildId uint64, userId uint64) bool {
	var node PermissionNode
	database.Database.Where(&PermissionNode{GuildId: guildId, UserId: userId}).Take(&node)
	return node.IsAdmin
}

func IsStaff(guildId uint64, userId uint64) bool {
	var node PermissionNode
	database.Database.Where(&PermissionNode{GuildId: guildId, UserId: userId}).Take(&node)
	return node.IsAdmin || node.IsSupport
}
