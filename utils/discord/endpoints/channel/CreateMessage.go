package channel

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils/discord"
)

type CreateMessageBody struct {
	Content string `json:"content"`
}

func CreateMessage(id int) discord.Endpoint {
	return discord.Endpoint{
		RequestType: discord.POST,
		AuthorizationType: discord.BOT,
		Endpoint: fmt.Sprintf("/channels/%d/messages", id),
	}
}
