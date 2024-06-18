package root

import (
	"context"
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/session"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rxdn/gdl/rest"
	"strconv"
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
	code, ok := ctx.GetQuery("code")
	if !ok {
		ctx.JSON(400, utils.ErrorStr("Discord provided invalid Oauth2 code"))
		return
	}

	res, err := discord.AccessToken(code)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Get ID + name
	currentUser, err := rest.GetCurrentUser(context.Background(), fmt.Sprintf("Bearer %s", res.AccessToken), nil)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	store := session.SessionData{
		AccessToken:  res.AccessToken,
		Expiry:       (time.Now().UnixNano() / int64(time.Second)) + int64(res.ExpiresIn),
		RefreshToken: res.RefreshToken,
		Name:         currentUser.Username,
		Avatar:       currentUser.AvatarUrl(256),
		HasGuilds:    false,
	}

	if err := utils.LoadGuilds(res.AccessToken, currentUser.Id); err == nil {
		store.HasGuilds = true
	} else {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":    strconv.FormatUint(currentUser.Id, 10),
		"timestamp": time.Now(),
	})

	str, err := token.SignedString([]byte(config.Conf.Server.Secret))
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if err := session.Store.Set(currentUser.Id, store); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"token":   str,
	})
}
