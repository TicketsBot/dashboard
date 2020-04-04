package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type CachedRole struct {
	AssociationId int    `gorm:"column:ASSOCIATIONID;primary_key;auto_increment"`
	GuildId       uint64 `gorm:"column:GUILDID"`
	UserId        uint64 `gorm:"column:USERID"`
	RoleId        uint64 `gorm:"column:ROLEID"`
}

func (CachedRole) TableName() string {
	return "cache_roles"
}

func DeleteRoles(guildId, userId uint64) {
	database.Database.Where(CachedRole{
		GuildId: guildId,
		UserId:  userId,
	}).Delete(CachedRole{})
}

// TODO: Cache invalidation
func CacheRole(guildId, userId, roleId uint64) {
	database.Database.Create(&CachedRole{
		GuildId: guildId,
		UserId:  userId,
		RoleId:  roleId,
	})
}

func GetCachedRoles(guildId, userId uint64, res chan []uint64) {
	var rows []CachedRole
	database.Database.Where(&CachedRole{
		GuildId: guildId,
		UserId:  userId,
	}).Find(&rows)

	roles := make([]uint64, 0)
	for _, row := range rows {
		roles = append(roles, row.RoleId)
	}

	res <- roles
}
