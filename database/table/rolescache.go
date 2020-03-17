package table

// Use an intermediary table to prevent a many-to-many relationship
type RoleCache struct {
	MemberId int   `gorm:"column:MEMBERID;primary_key"`
	RoleId   int64 `gorm:"column:ROLEID"`
}

func (RoleCache) TableName() string {
	return "cache_roles"
}
