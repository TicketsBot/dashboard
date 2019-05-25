package user

import "github.com/TicketsBot/GoPanel/utils/discord"

var CurrentUser = discord.Endpoint{
	RequestType: discord.GET,
	AuthorizationType: discord.BEARER,
	Endpoint: "/users/@me",
}
