package chatreplica

import (
	"fmt"
	"github.com/rxdn/gdl/objects/channel/message"
	"strconv"
)

func FromArchiveMessages(messages []message.Message, ticketId int) Payload {
	users := make(map[string]User)
	var wrappedMessages []Message // Cannot define length because of continue

	for _, msg := range messages {
		// If all 3 are missing, server will 400
		if msg.Content == "" && len(msg.Embeds) == 0 && len(msg.Attachments) == 0 {
			continue
		}

		wrappedMessages = append(wrappedMessages, Message{
			Id:          msg.Id,
			Type:        msg.Type,
			Author:      msg.Author.Id,
			Time:        msg.Timestamp.UnixMilli(),
			Content:     msg.Content,
			Embeds:      msg.Embeds,
			Attachments: msg.Attachments,
		})

		// Add user to entities map
		snowflake := strconv.FormatUint(msg.Author.Id, 10)
		if _, ok := users[snowflake]; !ok {
			var badge *Badge
			if msg.Author.Bot {
				badge = badgePtr(BadgeBot)
			}

			users[snowflake] = User{
				Avatar:        msg.Author.AvatarUrl(256),
				Username:      msg.Author.Username,
				Discriminator: fmt.Sprintf("%04d", msg.Author.Discriminator),
				Badge:         badge,
			}
		}
	}

	return Payload{
		Entities: Entities{
			Users:    users,
			Channels: make(map[string]Channel),
			Roles:    make(map[string]Role),
		},
		Messages:    wrappedMessages,
		ChannelName: fmt.Sprintf("ticket-%d", ticketId),
	}
}
