package chatreplica

import (
	"fmt"
	v2 "github.com/TicketsBot/logarchiver/model/v2"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/channel/message"
	"strconv"
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
		Avatar        string `json:"avatar"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Badge         *Badge `json:"badge,omitempty"`
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

// TODO: Use a generic ptr func
func badgePtr(b Badge) *Badge {
	return &b
}

func MessagesFromTranscript(messages []v2.Message) []Message {
	// Can't assign length as we might filter
	var wrappedMessages []Message

	for _, msg := range messages {
		if msg.Content == "" && len(msg.Embeds) == 0 && len(msg.Attachments) == 0 {
			continue
		}

		wrappedMessages = append(wrappedMessages, Message{
			Id:          msg.Id,
			Type:        message.MessageTypeDefault,
			Author:      msg.AuthorId,
			Time:        msg.Timestamp.UnixMilli(),
			Content:     msg.Content,
			Embeds:      msg.Embeds,
			Attachments: msg.Attachments,
		})
	}

	return wrappedMessages
}

func EntitiesFromTranscript(entities v2.Entities) Entities {
	users := make(map[string]User)
	for _, user := range entities.Users {
		var badge *Badge
		if user.Bot {
			badge = badgePtr(BadgeBot)
		}

		users[strconv.FormatUint(user.Id, 10)] = User{
			Avatar:        user.AvatarUrl(256),
			Username:      user.Username,
			Discriminator: fmt.Sprintf("%04d", user.Discriminator),
			Badge:         badge,
		}
	}

	channels := make(map[string]Channel)
	for _, channel := range entities.Channels {
		channels[strconv.FormatUint(channel.Id, 10)] = Channel{
			Name: channel.Name,
		}
	}

	roles := make(map[string]Role)
	for _, role := range entities.Roles {
		roles[strconv.FormatUint(role.Id, 10)] = Role{
			Name:  role.Name,
			Color: int(role.Colour),
		}
	}

	return Entities{
		Users:    users,
		Channels: channels,
		Roles:    roles,
	}
}
