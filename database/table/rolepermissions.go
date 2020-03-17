package table

import (
	"github.com/TicketsBot/GoPanel/database"
	"strconv"
)

type RolePermissions struct {
	GuildId int64 `gorm:"column:GUILDID"`
	RoleId  int64 `gorm:"column:ROLEID"`
	Support bool  `gorm:"column:ISSUPPORT"`
	Admin   bool  `gorm:"column:ISADMIN"`
}

func (RolePermissions) TableName() string {
	return "role_permissions"
}

func IsSupportRole(guild string, role string, ch chan bool) {
	guildId, err := strconv.ParseInt(guild, 10, 64); if err != nil {
		ch <- false
		return
	}

	rollId, err := strconv.ParseInt(role, 10, 64); if err != nil {
		ch <- false
		return
	}

	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: rollId}).First(&node)
	ch <- node.Support
}

func IsAdminRole(guild string, role string, ch chan bool) {
	guildId, err := strconv.ParseInt(guild, 10, 64); if err != nil {
		ch <- false
		return
	}

	rollId, err := strconv.ParseInt(role, 10, 64); if err != nil {
		ch <- false
		return
	}

	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: rollId}).First(&node)
	ch <- node.Admin
}

func GetAdminRoles(guild string, ch chan []int64) {
	guildId, err := strconv.ParseInt(guild, 10, 64); if err != nil {
		ch <- []int64{}
		return
	}

	var nodes []RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, Admin: true}).Find(&nodes)

	ids := make([]int64, 0)
	for _, node := range nodes {
		ids = append(ids, node.RoleId)
	}

	ch <- ids
}

func GetSupportRoles(guild string, ch chan []int64) {
	guildId, err := strconv.ParseInt(guild, 10, 64); if err != nil {
		ch <- []int64{}
		return
	}

	var nodes []RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, Support: true, Admin: false}).Find(&nodes)

	ids := make([]int64, 0)
	for _, node := range nodes {
		ids = append(ids, node.RoleId)
	}

	ch <- ids
}

func AddAdminRole(guild string, role string) {
	guildId, err := strconv.ParseInt(guild, 10, 64); if err != nil {
		return
	}

	roleId, err := strconv.ParseInt(role, 10, 64); if err != nil {
		return
	}

	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: roleId}).Assign(RolePermissions{Admin: true, Support: true}).FirstOrCreate(&node)
}

func AddSupportRole(guild string, role string) {
	guildId, err := strconv.ParseInt(guild, 10, 64); if err != nil {
		return
	}

	roleId, err := strconv.ParseInt(role, 10, 64); if err != nil {
		return
	}

	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: roleId}).Assign(RolePermissions{Support: true}).FirstOrCreate(&node)
}

func RemoveAdminRole(guild string, role string) {
	guildId, err := strconv.ParseInt(guild, 10, 64); if err != nil {
		return
	}

	roleId, err := strconv.ParseInt(role, 10, 64); if err != nil {
		return
	}

	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: roleId}).Take(&node)
	database.Database.Model(&node).Where("GUILDID = ? AND ROLEID = ?", guildId, roleId).Update("ISADMIN", false)
}

func RemoveSupportRole(guild string, role string) {
	guildId, err := strconv.ParseInt(guild, 10, 64); if err != nil {
		return
	}

	roleId, err := strconv.ParseInt(role, 10, 64); if err != nil {
		return
	}

	var node RolePermissions
	database.Database.Where(RolePermissions{GuildId: guildId, RoleId: roleId}).Take(&node)
	database.Database.Model(&node).Where("GUILDID = ? AND ROLEID = ?", guildId, roleId).Updates(map[string]interface{}{
		"ISADMIN": false,
		"ISSUPPORT": false,
	})
}
