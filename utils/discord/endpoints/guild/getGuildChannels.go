package guild

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils/discord"
)

func GetGuildChannels(id int) discord.Endpoint {
	return discord.Endpoint{
		RequestType: discord.GET,
		AuthorizationType: discord.BOT,
		Endpoint: fmt.Sprintf("/guilds/%d/channels", id),
	}
}
