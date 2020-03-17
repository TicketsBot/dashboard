package guild

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils/discord"
	"strconv"
)

func GetGuildMember(guildId, userId int) discord.Endpoint {
	return discord.Endpoint{
		RequestType:       discord.GET,
		AuthorizationType: discord.BOT,
		Endpoint:          fmt.Sprintf("/guilds/%s/members/%s", strconv.Itoa(guildId), strconv.Itoa(userId)),
	}
}
