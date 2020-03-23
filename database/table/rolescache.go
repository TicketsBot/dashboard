package table

import "github.com/TicketsBot/TicketsGo/database"

type CachedRole struct {
	AssociationId int   `gorm:"column:ASSOCIATIONID;primary_key;auto_increment"`
	GuildId       int64 `gorm:"column:GUILDID"`
	UserId        int64 `gorm:"column:USERID"`
	RoleId        int64 `gorm:"column:ROLEID"`
}

func (CachedRole) TableName() string {
	return "cache_roles"
}

func DeleteRoles(guildId, userId int64) {
	database.Db.Where(CachedRole{
		GuildId: guildId,
		UserId:  userId,
	}).Delete(CachedRole{})
}

// TODO: Cache invalidation
func CacheRole(guildId, userId, roleId int64) {
	database.Db.Create(&CachedRole{
		GuildId: guildId,
		UserId:  userId,
		RoleId:  roleId,
	})
}

func GetCachedRoles(guildId, userId int64, res chan []int64) {
	var rows []CachedRole
	database.Db.Where(&CachedRole{
		GuildId: guildId,
		UserId:  userId,
	}).Find(&rows)

	roles := make([]int64, 0)
	for _, row := range rows {
		roles = append(roles, row.RoleId)
	}
	res <- roles
}
