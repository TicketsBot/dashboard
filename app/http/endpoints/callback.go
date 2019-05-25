package endpoints

import (
	"github.com/TicketsBot/GoPanel/utils/discord"
	"github.com/TicketsBot/GoPanel/utils/discord/endpoints/user"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

type(
	TokenData struct {
		ClientId string `qs:"client_id"`
		ClientSecret string `qs:"client_secret"`
		GrantType string `qs:"grant_type"`
		Code string `qs:"code"`
		RedirectUri string `qs:"redirect_uri"`
		Scope string `qs:"scope"`
	}

	TokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType string `json:"token_type"`
		ExpiresIn int `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope string `json:"scope"`
	}
)

func CallbackHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	code := ctx.DefaultQuery("code", "")
	if code == "" {
		ctx.String(400, "Discord provided an invalid Oauth2 code")
		return
	}

	res, err := discord.AccessToken(code); if err != nil {
		ctx.String(500, err.Error())
	}

	store.Set("access_token", res.AccessToken)
	store.Set("refresh_token", res.RefreshToken)
	store.Set("expiry", (time.Now().UnixNano() / int64(time.Second)) + int64(res.ExpiresIn))

	// Get ID + name
	var currentUser objects.User
	err = user.CurrentUser.Request(store, nil, nil, &currentUser); if err != nil {
		ctx.String(500, err.Error())
		return
	}

	store.Set("userid", currentUser.Id)
	store.Set("name", currentUser.Username)

	// Get Guilds
	var currentUserGuilds []objects.Guild
	err = user.CurrentUserGuilds.Request(store, nil, nil, &currentUserGuilds); if err != nil {
		ctx.String(500, err.Error())
	}
}
