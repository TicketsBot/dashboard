package channel

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils/discord"
)

func GetChannelMessages(id int) discord.Endpoint {
	return discord.Endpoint{
		RequestType: discord.GET,
		AuthorizationType: discord.BOT,
		Endpoint: fmt.Sprintf("/channels/%d/messages", id),
	}
}
