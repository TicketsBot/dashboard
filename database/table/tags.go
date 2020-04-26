package table

import (
	"github.com/TicketsBot/GoPanel/database"
	uuid "github.com/satori/go.uuid"
)

type Tag struct {
	Uuid string `gorm:"column:UUID;type:varchar(36);unique;primary_key"`
	Id string `gorm:"column:ID;type:varchar(16)"`
	Guild uint64 `gorm:"column:GUILDID"`
	Content string `gorm:"column:TEXT;type:TEXT"`
}

func (Tag) TableName() string {
	return "cannedresponses"
}

func GetTag(guild uint64, id string, ch chan string) {
	var node Tag
	database.Database.Where(Tag{Id: id, Guild: guild}).Take(&node)
	ch <- node.Content
}

func GetTags(guild uint64, ch chan []Tag) {
	var rows []Tag
	database.Database.Where(Tag{Guild: guild}).Find(&rows)
	ch <- rows
}

func AddTag(guild uint64, id string, content string) error {
	return database.Database.Create(&Tag{
		Uuid: uuid.NewV4().String(),
		Id: id,
		Guild: guild,
		Content: content,
	}).Error
}

func DeleteTag(guild uint64, id string) error {
	return database.Database.Where(Tag{Id: id, Guild: guild}).Delete(&Tag{}).Error
}
