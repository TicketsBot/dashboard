package table

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
)

type UsernameNode struct {
	Id int64  `gorm:"column:USERID;primary_key"`
	Name string `gorm:"column:USERNAME;type:text"` // Base 64 encoded
	Discriminator string `gorm:"column:DISCRIM;type:varchar(4)"`
	Avatar string `gorm:"column:AVATARHASH;type:varchar(100)"`
}

func (UsernameNode) TableName() string {
	return "usernames"
}

func GetUsername(id int64) string {
	node := UsernameNode{Name: "Unknown"}
	database.Database.Where(&UsernameNode{Id: id}).First(&node)
	return utils.Base64Decode(node.Name)
}

func GetUserNodes(ids []int64) []UsernameNode {
	var nodes []UsernameNode
	database.Database.Where(ids).Find(&nodes)
	return nodes
}

func GetUserId(name, discrim string) int64 {
	var node UsernameNode
	database.Database.Where(&UsernameNode{Name: utils.Base64Encode(name), Discriminator: discrim}).First(&node)
	return node.Id
}
