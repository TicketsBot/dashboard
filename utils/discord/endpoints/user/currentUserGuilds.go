package user

import "github.com/TicketsBot/GoPanel/utils/discord"

var CurrentUserGuilds = discord.Endpoint{
	RequestType: discord.GET,
	Endpoint: "/users/@me/guilds",
}
