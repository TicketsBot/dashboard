package user

import "github.com/TicketsBot/GoPanel/utils/discord"

var CurrentUser = discord.Endpoint{
	RequestType: discord.GET,
	Endpoint: "/users/@me",
}
