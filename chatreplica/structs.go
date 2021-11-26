package chatreplica

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/user"
)

type (
	Payload struct {
		Entities    Entities  `json:"entities"`
		Messages    []Message `json:"messages"`
		ChannelName string    `json:"channel_name"`
	}

	// Entities Snowflake -> Entity map
	Entities struct {
		Users    map[string]User    `json:"users"`
		Channels map[string]Channel `json:"channels"`
		Roles    map[string]Role    `json:"roles"`
	}

	User struct {
		Avatar        string             `json:"avatar"`
		Username      string             `json:"username"`
		Discriminator user.Discriminator `json:"discriminator"`
		Badge         *Badge             `json:"badge,omitempty"`
	}

	Channel struct {
		Name string `json:"name"`
	}

	Role struct {
		Name  string `json:"name"`
		Color int    `json:"color"`
	}

	Message struct {
		Id          uint64               `json:"id,string"`
		Type        message.MessageType  `json:"type"`
		Author      uint64               `json:"author,string"`
		Time        int64                `json:"time"` // Unix seconds
		Content     string               `json:"content"`
		Embeds      []embed.Embed        `json:"embeds,omitempty"`
		Attachments []channel.Attachment `json:"attachments,omitempty"`
	}
)

type Badge string

const (
	BadgeBot Badge = "bot"
)

func badgePtr(b Badge) *Badge {
    return &b
}
