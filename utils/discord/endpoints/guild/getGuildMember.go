package guild

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils/discord"
)

func GetGuildMember(guildId, userId uint64) discord.Endpoint {
	return discord.Endpoint{
		RequestType:       discord.GET,
		AuthorizationType: discord.BOT,
		Endpoint:          fmt.Sprintf("/guilds/%d/members/%d", guildId, userId),
	}
}
