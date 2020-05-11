package root

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord"
	userEndpoint "github.com/TicketsBot/GoPanel/utils/discord/endpoints/user"
	"github.com/TicketsBot/database"
	"github.com/apex/log"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/user"
	"strings"
	"time"
)

type (
	TokenData struct {
		ClientId     string `qs:"client_id"`
		ClientSecret string `qs:"client_secret"`
		GrantType    string `qs:"grant_type"`
		Code         string `qs:"code"`
		RedirectUri  string `qs:"redirect_uri"`
		Scope        string `qs:"scope"`
	}

	TokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}
)

func CallbackHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	if utils.IsLoggedIn(store) {
		ctx.Redirect(302, config.Conf.Server.BaseUrl)
		return
	}

	code, ok := ctx.GetQuery("code")
	if !ok {
		ctx.String(400, "Discord provided an invalid Oauth2 code")
		return
	}

	res, err := discord.AccessToken(code)
	if err != nil {
		ctx.String(500, err.Error())
	}

	store.Set("access_token", res.AccessToken)
	store.Set("refresh_token", res.RefreshToken)
	store.Set("expiry", (time.Now().UnixNano()/int64(time.Second))+int64(res.ExpiresIn))

	// Get ID + name
	var currentUser user.User
	err, _ = userEndpoint.CurrentUser.Request(store, nil, nil, &currentUser)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	store.Set("csrf", utils.RandStringRunes(32))

	store.Set("userid", currentUser.Id)
	store.Set("name", currentUser.Username)
	store.Set("avatar", fmt.Sprintf("https://cdn.discordapp.com/avatars/%d/%s.webp", currentUser.Id, currentUser.Avatar))
	if err = store.Save(); err != nil {
		log.Error(err.Error())
	}

	handleRedirect(ctx)

	// Cache guilds because Discord takes like 2 whole seconds to return then
	go func() {
		var guilds []*guild.Guild
		err, _ = userEndpoint.CurrentUserGuilds.Request(store, nil, nil, &guilds)
		if err != nil {
			log.Error(err.Error())
			return
		}

		var wrappedGuilds []database.UserGuild

		// endpoint's partial guild doesn't include ownerid
		// we only user cached guilds on the index page, so it doesn't matter if we don't have have the real owner id
		// if the user isn't the owner, as we pull from the cache on other endpoints
		for _, guild := range guilds {
			if guild.Owner {
				guild.OwnerId = currentUser.Id
			}

			wrappedGuilds = append(wrappedGuilds, database.UserGuild{
				GuildId:         guild.Id,
				Name:            guild.Name,
				Owner:           guild.Owner,
				UserPermissions: int32(guild.Permissions),
			})
		}

		// TODO: Error handling
		go dbclient.Client.UserGuilds.Set(currentUser.Id, wrappedGuilds)
	}()
}

func handleRedirect(ctx *gin.Context) {
	state := strings.Split(ctx.Query("state"), ".")

	if len(state) == 3 && state[0] == "viewlog" {
		ctx.Redirect(302, fmt.Sprintf("%s/manage/%s/logs/view/%s", config.Conf.Server.BaseUrl, state[1], state[2]))
	} else {
		ctx.Redirect(302, config.Conf.Server.BaseUrl)
	}
}
