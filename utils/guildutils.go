package utils

import (
	"fmt"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/database"
	"github.com/rxdn/gdl/rest"
)

func LoadGuilds(accessToken string, userId uint64) error {
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	data := rest.CurrentUserGuildsData{
		Limit: 100,
	}

	guilds, err := rest.GetCurrentUserGuilds(authHeader, nil, data)
	if err != nil {
		return err
	}

	var wrappedGuilds []database.UserGuild

	// endpoint's partial guild doesn't includes ownerid
	// we only user cached guilds on the index page, so it doesn't matter if we don't have have the real owner id
	// if the user isn't the owner, as we pull from the cache on other endpoints
	for _, guild := range guilds {
		wrappedGuilds = append(wrappedGuilds, database.UserGuild{
			GuildId:         guild.Id,
			Name:            guild.Name,
			Owner:           guild.Owner,
			UserPermissions: int32(guild.Permissions),
			Icon:            guild.Icon,
		})
	}

	return dbclient.Client.UserGuilds.Set(userId, wrappedGuilds)
}
