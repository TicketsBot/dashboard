package root

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord"
	"github.com/apex/log"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
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

	if utils.IsLoggedIn(store) && store.Get("has_guilds") == true {
		ctx.Redirect(302, config.Conf.Server.BaseUrl)
		return
	}

	code, ok := ctx.GetQuery("code")
	if !ok {
		utils.ErrorPage(ctx, 400, "Discord provided invalid Oauth2 code")
		return
	}

	res, err := discord.AccessToken(code)
	if err != nil {
		utils.ErrorPage(ctx, 500, err.Error())
		return
	}

	store.Set("access_token", res.AccessToken)
	store.Set("refresh_token", res.RefreshToken)
	store.Set("expiry", (time.Now().UnixNano()/int64(time.Second))+int64(res.ExpiresIn))

	// Get ID + name
	currentUser, err := rest.GetCurrentUser(fmt.Sprintf("Bearer %s", res.AccessToken), nil)
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	store.Set("csrf", utils.RandString(32))

	store.Set("userid", currentUser.Id)
	store.Set("name", currentUser.Username)
	store.Set("avatar", currentUser.AvatarUrl(256))
	store.Save()

	if err := utils.LoadGuilds(res.AccessToken, currentUser.Id); err == nil {
		store.Set("has_guilds", true)
		store.Save()
	} else {
		log.Error(err.Error())
	}

	handleRedirect(ctx)
}

func handleRedirect(ctx *gin.Context) {
	state := strings.Split(ctx.Query("state"), ".")

	if len(state) == 3 && state[0] == "viewlog" {
		ctx.Redirect(302, fmt.Sprintf("%s/manage/%s/logs/view/%s", config.Conf.Server.BaseUrl, state[1], state[2]))
	} else {
		ctx.Redirect(302, config.Conf.Server.BaseUrl)
	}
}
