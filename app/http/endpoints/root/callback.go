package root

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord"
	userEndpoint "github.com/TicketsBot/GoPanel/utils/discord/endpoints/user"
	"github.com/apex/log"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/user"
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

	code := ctx.DefaultQuery("code", "")
	if code == "" {
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
	store.Set("avatar", fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.webp", currentUser.Id, currentUser.Avatar))
	if err = store.Save(); err != nil {
		log.Error(err.Error())
	}

	ctx.Redirect(302, config.Conf.Server.BaseUrl)

	// Cache guilds because Discord takes like 2 whole seconds to return then
	go func() {
		var guilds []guild.Guild
		err, _ = userEndpoint.CurrentUserGuilds.Request(store, nil, nil, &guilds)
		if err != nil {
			log.Error(err.Error())
			return
		}

		marshalled, err := json.Marshal(guilds)
		if err != nil {
			log.Error(err.Error())
			return
		}

		// TODO: unfuck this
		table.UpdateGuilds(currentUser.Id, base64.StdEncoding.EncodeToString(marshalled))
	}()
}