package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type RolePermissions struct {
	GuildId uint64 `gorm:"column:GUILDID"`
	RoleId  uint64 `gorm:"column:ROLEID"`
	Support bool   `gorm:"column:ISSUPPORT"`
	Admin   bool   `gorm:"column:ISADMIN"`
}

func (RolePermissions) TableName() string {
	return "role_permissions"
}

func IsSupportRole(guildId, roleId uint64, ch chan bool) {
	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: roleId}).First(&node)
	ch <- node.Support
}

func IsAdminRole(guildId, roleId uint64, ch chan bool) {
	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: roleId}).First(&node)
	ch <- node.Admin
}

func GetAdminRoles(guildId uint64, ch chan []uint64) {
	var nodes []RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, Admin: true}).Find(&nodes)

	ids := make([]uint64, 0)
	for _, node := range nodes {
		ids = append(ids, node.RoleId)
	}

	ch <- ids
}

func GetSupportRoles(guildId uint64, ch chan []uint64) {
	var nodes []RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, Support: true, Admin: false}).Find(&nodes)

	ids := make([]uint64, 0)
	for _, node := range nodes {
		ids = append(ids, node.RoleId)
	}

	ch <- ids
}

func AddAdminRole(guildId, roleId uint64) {
	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: roleId}).Assign(RolePermissions{Admin: true, Support: true}).FirstOrCreate(&node)
}

func AddSupportRole(guildId, roleId uint64) {
	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: roleId}).Assign(RolePermissions{Support: true}).FirstOrCreate(&node)
}

func RemoveAdminRole(guildId, roleId uint64) {
	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: roleId}).Take(&node)
	database.Database.Model(&node).Where("GUILDID = ? AND ROLEID = ?", guildId, roleId).Update("ISADMIN", false)
}

func RemoveSupportRole(guildId, roleId uint64) {
	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: roleId}).Take(&node)
	database.Database.Model(&node).Where("GUILDID = ? AND ROLEID = ?", guildId, roleId).Updates(map[string]interface{}{
		"ISADMIN":   false,
		"ISSUPPORT": false,
	})
}
