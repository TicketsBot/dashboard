package api

type entityType int

const (
	entityTypeUser entityType = iota
	entityTypeRole
)

var entityTypes = map[int]entityType{
	int(entityTypeUser): entityTypeUser,
	int(entityTypeRole): entityTypeRole,
}

type entity struct {
	Id   uint64     `json:"id,string"`
	Name string     `json:"name"`
	Type entityType `json:"type"`
}
