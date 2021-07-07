package http

import (
	"github.com/TicketsBot/GoPanel/app/http/endpoints/api"
	api_autoclose "github.com/TicketsBot/GoPanel/app/http/endpoints/api/autoclose"
	api_blacklist "github.com/TicketsBot/GoPanel/app/http/endpoints/api/blacklist"
	api_panels "github.com/TicketsBot/GoPanel/app/http/endpoints/api/panel"
	api_settings "github.com/TicketsBot/GoPanel/app/http/endpoints/api/settings"
	api_tags "github.com/TicketsBot/GoPanel/app/http/endpoints/api/tags"
	api_team "github.com/TicketsBot/GoPanel/app/http/endpoints/api/team"
	api_ticket "github.com/TicketsBot/GoPanel/app/http/endpoints/api/ticket"
	api_transcripts "github.com/TicketsBot/GoPanel/app/http/endpoints/api/transcripts"
	api_whitelabel "github.com/TicketsBot/GoPanel/app/http/endpoints/api/whitelabel"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/root"
	"github.com/TicketsBot/GoPanel/app/http/middleware"
	"github.com/TicketsBot/GoPanel/app/http/session"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"log"
	"time"
)

func StartServer() {
	log.Println("Starting HTTP server")

	router := gin.Default()

	// Sessions
	session.Store = session.NewRedisStore()

	// Handle static asset requests
	router.Use(static.Serve("/assets/", static.LocalFile("./public/static", false)))

	router.Use(gin.Recovery())
	router.Use(createLimiter(600, time.Minute*10))

	router.Use(middleware.Cors(config.Conf))

	router.GET("/webchat", root.WebChatWs)

	router.POST("/callback", middleware.VerifyXTicketsHeader, root.CallbackHandler)
	router.POST("/logout", middleware.VerifyXTicketsHeader, middleware.AuthenticateToken, root.LogoutHandler)

	apiGroup := router.Group("/api", middleware.VerifyXTicketsHeader, middleware.AuthenticateToken)
	{
		apiGroup.GET("/session", api.SessionHandler)
	}

	guildAuthApiAdmin := apiGroup.Group("/:id", middleware.AuthenticateGuild(true, permission.Admin))
	guildAuthApiSupport := apiGroup.Group("/:id", middleware.AuthenticateGuild(true, permission.Support))
	guildApiNoAuth := apiGroup.Group("/:id", middleware.ParseGuildId)
	{
		guildAuthApiSupport.GET("/channels", api.ChannelsHandler)
		guildAuthApiSupport.GET("/premium", api.PremiumHandler)
		guildAuthApiSupport.GET("/user/:user", api.UserHandler)
		guildAuthApiSupport.GET("/roles", api.RolesHandler)
		guildAuthApiSupport.GET("/members/search", createLimiter(5, time.Second), createLimiter(10, time.Second * 30), createLimiter(75, time.Minute * 30), api.SearchMembers)

		guildAuthApiAdmin.GET("/settings", api_settings.GetSettingsHandler)
		guildAuthApiAdmin.POST("/settings", api_settings.UpdateSettingsHandler)

		guildAuthApiSupport.GET("/blacklist", api_blacklist.GetBlacklistHandler)
		guildAuthApiSupport.POST("/blacklist/:user", api_blacklist.AddBlacklistHandler)
		guildAuthApiSupport.DELETE("/blacklist/:user", api_blacklist.RemoveBlacklistHandler)

		guildAuthApiAdmin.GET("/panels", api_panels.ListPanels)
		guildAuthApiAdmin.POST("/panels", api_panels.CreatePanel)
		guildAuthApiAdmin.PATCH("/panels/:panelid", api_panels.UpdatePanel)
		guildAuthApiAdmin.DELETE("/panels/:panelid", api_panels.DeletePanel)

		guildAuthApiAdmin.GET("/multipanels", api_panels.MultiPanelList)
		guildAuthApiAdmin.POST("/multipanels", api_panels.MultiPanelCreate)
		guildAuthApiAdmin.PATCH("/multipanels/:panelid", api_panels.MultiPanelUpdate)
		guildAuthApiAdmin.DELETE("/multipanels/:panelid", api_panels.MultiPanelDelete)

		guildAuthApiSupport.GET("/transcripts", createLimiter(5, 5 * time.Second), createLimiter(20, time.Minute), api_transcripts.ListTranscripts)
		// Allow regular users to get their own transcripts, make sure you check perms inside
		guildApiNoAuth.GET("/transcripts/:ticketId", createLimiter(10, 10 * time.Second), api_transcripts.GetTranscriptHandler)

		guildAuthApiSupport.GET("/tickets", api_ticket.GetTickets)
		guildAuthApiSupport.GET("/tickets/:ticketId", api_ticket.GetTicket)
		guildAuthApiSupport.POST("/tickets/:ticketId", createLimiter(5, time.Second * 5), api_ticket.SendMessage)
		guildAuthApiSupport.DELETE("/tickets/:ticketId", api_ticket.CloseTicket)

		guildAuthApiSupport.GET("/tags", api_tags.TagsListHandler)
		guildAuthApiSupport.PUT("/tags", api_tags.CreateTag)
		guildAuthApiSupport.DELETE("/tags/:tag", api_tags.DeleteTag)

		guildAuthApiAdmin.GET("/claimsettings", api_settings.GetClaimSettings)
		guildAuthApiAdmin.POST("/claimsettings", api_settings.PostClaimSettings)

		guildAuthApiAdmin.GET("/autoclose", api_autoclose.GetAutoClose)
		guildAuthApiAdmin.POST("/autoclose", api_autoclose.PostAutoClose)

		guildAuthApiAdmin.GET("/team", api_team.GetTeams)
		guildAuthApiAdmin.GET("/team/:teamid", createLimiter(10, time.Second * 30), api_team.GetMembers)
		guildAuthApiAdmin.POST("/team", createLimiter(10, time.Minute), api_team.CreateTeam)
		guildAuthApiAdmin.PUT("/team/:teamid/:snowflake", createLimiter(5, time.Second * 10), api_team.AddMember)
		guildAuthApiAdmin.DELETE("/team/:teamid", api_team.DeleteTeam)
		guildAuthApiAdmin.DELETE("/team/:teamid/:snowflake", createLimiter(30, time.Minute), api_team.RemoveMember)
	}

	userGroup := router.Group("/user", middleware.AuthenticateToken)
	{
		userGroup.GET("/guilds", api.GetGuilds)
		userGroup.POST("/guilds/reload", api.ReloadGuildsHandler)
		userGroup.GET("/permissionlevel", api.GetPermissionLevel)

		{
			whitelabelGroup := userGroup.Group("/whitelabel", middleware.VerifyWhitelabel(true))

			whitelabelGroup.GET("/", api_whitelabel.WhitelabelGet)
			whitelabelGroup.GET("/errors", api_whitelabel.WhitelabelGetErrors)
			whitelabelGroup.GET("/guilds", api_whitelabel.WhitelabelGetGuilds)
			whitelabelGroup.GET("/public-key", api_whitelabel.WhitelabelGetPublicKey)
			whitelabelGroup.POST("/public-key", api_whitelabel.WhitelabelPostPublicKey)
			whitelabelGroup.POST("/create-interactions", api_whitelabel.GetWhitelabelCreateInteractions())

			whitelabelGroup.POST("/", createLimiter(10, time.Minute), api_whitelabel.WhitelabelPost)
			whitelabelGroup.POST("/status", createLimiter(1, time.Second*5), api_whitelabel.WhitelabelStatusPost)
		}
	}

	if err := router.Run(config.Conf.Server.Host); err != nil {
		panic(err)
	}
}

func serveTemplate(templateName string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		guildId := ctx.Keys["guildid"].(uint64)
		userId := ctx.Keys["userid"].(uint64)

		store, err := session.Store.Get(userId)
		if err != nil {
			if err == session.ErrNoSession {
				ctx.JSON(401, gin.H{
					"success": false,
					"auth": true,
				})
			} else {
				ctx.JSON(500, utils.ErrorJson(err))
			}

			return
		}

		ctx.HTML(200, templateName, gin.H{
			"name":         store.Name,
			"guildId":      guildId,
			"avatar":       store.Avatar,
			"baseUrl":      config.Conf.Server.BaseUrl,
		})
	}
}

func createLimiter(limit int64, period time.Duration) func(*gin.Context) {
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}

	return mgin.NewMiddleware(limiter.New(store, rate))
}
