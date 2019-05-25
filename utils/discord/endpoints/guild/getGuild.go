package guild

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils/discord"
)

func GetGuild(id int) discord.Endpoint {
	return discord.Endpoint{
		RequestType: discord.GET,
		Endpoint: fmt.Sprintf("/guilds/%d", id),
	}
}
