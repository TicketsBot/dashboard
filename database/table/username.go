package table

import (
	"encoding/base64"
	"github.com/TicketsBot/GoPanel/database"
)

type UsernameNode struct {
	Id int64  `gorm:"column:USERID;primary_key"`
	Name string `gorm:"column:USERNAME;type:text"` // Base 64 encoded
}

func (UsernameNode) TableName() string {
	return "usernames"
}

func GetUsername(id int64) string {
	node := UsernameNode{Name: "Unknown"}
	database.Database.Where(&UsernameNode{Id: id}).First(&node)
	return base64Decode(node.Name)
}

func GetUsernames(ids []int64) map[int64]string {
	var nodes []UsernameNode
	database.Database.Where(ids).Find(&nodes)

	m := make(map[int64]string)
	for _, node := range nodes {
		m[node.Id] = base64Decode(node.Name)
	}

	return m
}

func base64Decode(s string) string {
	b, err := base64.StdEncoding.DecodeString(s); if err != nil {
		return ""
	}
	return string(b)
}
