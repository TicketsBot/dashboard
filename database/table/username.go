package table

import (
	"github.com/TicketsBot/GoPanel/database"
)

type UsernameNode struct {
	Id            uint64 `gorm:"column:USERID;primary_key"`
	Name          string `gorm:"column:USERNAME;type:text"` // Base 64 encoded
	Discriminator string `gorm:"column:DISCRIM;type:varchar(4)"`
	Avatar        string `gorm:"column:AVATARHASH;type:varchar(100)"`
}

func (UsernameNode) TableName() string {
	return "usernames"
}

func GetUsername(id uint64, ch chan string) {
	node := UsernameNode{Name: "Unknown"}
	database.Database.Where(&UsernameNode{Id: id}).First(&node)
	ch <- node.Name
}

func GetUserNodes(ids []uint64) []UsernameNode {
	var nodes []UsernameNode
	database.Database.Where(ids).Find(&nodes)
	return nodes
}

func GetUserId(name, discrim string) uint64 {
	var node UsernameNode
	database.Database.Where(&UsernameNode{Name: name, Discriminator: discrim}).First(&node)
	return node.Id
}
