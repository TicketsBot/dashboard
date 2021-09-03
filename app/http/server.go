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
	"github.com/TicketsBot/common/permission"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func StartServer() {
	log.Println("Starting HTTP server")

	router := gin.Default()

	router.RemoteIPHeaders = config.Conf.Server.RealIpHeaders
	if err := router.SetTrustedProxies(config.Conf.Server.TrustedProxies); err != nil {
		panic(err)
	}

	// Sessions
	session.Store = session.NewRedisStore()

	// Handle static asset requests
	router.Use(static.Serve("/assets/", static.LocalFile("./public/static", false)))

	router.Use(gin.Recovery())
	router.Use(middleware.MultiReadBody)
	router.Use(sentrygin.New(sentrygin.Options{})) // Defaults are ok

	router.Use(rl(middleware.RateLimitTypeIp, 60, time.Minute))
	router.Use(rl(middleware.RateLimitTypeUser, 60, time.Minute))
	router.Use(rl(middleware.RateLimitTypeGuild, 600, time.Minute*5))

	router.Use(middleware.Cors(config.Conf))

	router.GET("/webchat", root.WebChatWs, middleware.Logging)

	router.POST("/callback", middleware.VerifyXTicketsHeader, root.CallbackHandler, middleware.Logging)
	router.POST("/logout", middleware.VerifyXTicketsHeader, middleware.AuthenticateToken, root.LogoutHandler, middleware.Logging)

	apiGroup := router.Group("/api", middleware.VerifyXTicketsHeader, middleware.AuthenticateToken)
	{
		apiGroup.GET("/session", api.SessionHandler, middleware.Logging)
	}

	guildAuthApiAdmin := apiGroup.Group("/:id", middleware.AuthenticateGuild(true, permission.Admin))
	guildAuthApiSupport := apiGroup.Group("/:id", middleware.AuthenticateGuild(true, permission.Support))
	guildApiNoAuth := apiGroup.Group("/:id", middleware.ParseGuildId)
	{
		guildAuthApiSupport.GET("/channels", api.ChannelsHandler, middleware.Logging)
		guildAuthApiSupport.GET("/premium", api.PremiumHandler, middleware.Logging)
		guildAuthApiSupport.GET("/user/:user", api.UserHandler, middleware.Logging)
		guildAuthApiSupport.GET("/roles", api.RolesHandler, middleware.Logging)
		guildAuthApiSupport.GET("/members/search",
			rl(middleware.RateLimitTypeGuild, 5, time.Second),
			rl(middleware.RateLimitTypeGuild, 10, time.Second*30),
			rl(middleware.RateLimitTypeGuild, 75, time.Minute*30),
			api.SearchMembers,
			middleware.Logging,
		)

		guildAuthApiAdmin.GET("/settings", api_settings.GetSettingsHandler, middleware.Logging)
		guildAuthApiAdmin.POST("/settings", api_settings.UpdateSettingsHandler, middleware.Logging)

		guildAuthApiSupport.GET("/blacklist", api_blacklist.GetBlacklistHandler, middleware.Logging)
		guildAuthApiSupport.POST("/blacklist/:user", api_blacklist.AddBlacklistHandler, middleware.Logging)
		guildAuthApiSupport.DELETE("/blacklist/:user", api_blacklist.RemoveBlacklistHandler, middleware.Logging)

		guildAuthApiAdmin.GET("/panels", api_panels.ListPanels, middleware.Logging)
		guildAuthApiAdmin.POST("/panels", api_panels.CreatePanel, middleware.Logging)
		guildAuthApiAdmin.PATCH("/panels/:panelid", api_panels.UpdatePanel, middleware.Logging)
		guildAuthApiAdmin.DELETE("/panels/:panelid", api_panels.DeletePanel, middleware.Logging)

		guildAuthApiAdmin.GET("/multipanels", api_panels.MultiPanelList, middleware.Logging)
		guildAuthApiAdmin.POST("/multipanels", api_panels.MultiPanelCreate, middleware.Logging)
		guildAuthApiAdmin.PATCH("/multipanels/:panelid", api_panels.MultiPanelUpdate, middleware.Logging)
		guildAuthApiAdmin.DELETE("/multipanels/:panelid", api_panels.MultiPanelDelete, middleware.Logging)

		// Should be a GET, but easier to take a body for development purposes
		guildAuthApiSupport.POST("/transcripts",
			rl(middleware.RateLimitTypeUser, 5, 5*time.Second),
			rl(middleware.RateLimitTypeUser, 20, time.Minute),
			api_transcripts.ListTranscripts,
			middleware.Logging,
		)

		// Allow regular users to get their own transcripts, make sure you check perms inside
		guildApiNoAuth.GET("/transcripts/:ticketId", rl(middleware.RateLimitTypeGuild, 10, 10*time.Second), api_transcripts.GetTranscriptHandler, middleware.Logging)

		guildAuthApiSupport.GET("/tickets", api_ticket.GetTickets, middleware.Logging)
		guildAuthApiSupport.GET("/tickets/:ticketId", api_ticket.GetTicket, middleware.Logging)
		guildAuthApiSupport.POST("/tickets/:ticketId", rl(middleware.RateLimitTypeGuild, 5, time.Second*5), api_ticket.SendMessage, middleware.Logging)
		guildAuthApiSupport.DELETE("/tickets/:ticketId", api_ticket.CloseTicket, middleware.Logging)

		guildAuthApiSupport.GET("/tags", api_tags.TagsListHandler, middleware.Logging)
		guildAuthApiSupport.PUT("/tags", api_tags.CreateTag, middleware.Logging)
		guildAuthApiSupport.DELETE("/tags", api_tags.DeleteTag, middleware.Logging)

		guildAuthApiAdmin.GET("/claimsettings", api_settings.GetClaimSettings, middleware.Logging)
		guildAuthApiAdmin.POST("/claimsettings", api_settings.PostClaimSettings, middleware.Logging)

		guildAuthApiAdmin.GET("/autoclose", api_autoclose.GetAutoClose, middleware.Logging)
		guildAuthApiAdmin.POST("/autoclose", api_autoclose.PostAutoClose, middleware.Logging)

		guildAuthApiAdmin.GET("/team", api_team.GetTeams, middleware.Logging)
		guildAuthApiAdmin.GET("/team/:teamid", rl(middleware.RateLimitTypeUser ,10, time.Second*30), api_team.GetMembers, middleware.Logging)
		guildAuthApiAdmin.POST("/team", rl(middleware.RateLimitTypeUser, 10, time.Minute), api_team.CreateTeam, middleware.Logging)
		guildAuthApiAdmin.PUT("/team/:teamid/:snowflake", rl(middleware.RateLimitTypeGuild, 5, time.Second*10), api_team.AddMember, middleware.Logging)
		guildAuthApiAdmin.DELETE("/team/:teamid", api_team.DeleteTeam, middleware.Logging)
		guildAuthApiAdmin.DELETE("/team/:teamid/:snowflake", rl(middleware.RateLimitTypeGuild, 30, time.Minute), api_team.RemoveMember, middleware.Logging)
	}

	userGroup := router.Group("/user", middleware.AuthenticateToken)
	{
		userGroup.GET("/guilds", api.GetGuilds, middleware.Logging)
		userGroup.POST("/guilds/reload", api.ReloadGuildsHandler, middleware.Logging)
		userGroup.GET("/permissionlevel", api.GetPermissionLevel, middleware.Logging)

		{
			whitelabelGroup := userGroup.Group("/whitelabel", middleware.VerifyWhitelabel(true))

			whitelabelGroup.GET("/", api_whitelabel.WhitelabelGet, middleware.Logging)
			whitelabelGroup.GET("/errors", api_whitelabel.WhitelabelGetErrors, middleware.Logging)
			whitelabelGroup.GET("/guilds", api_whitelabel.WhitelabelGetGuilds, middleware.Logging)
			whitelabelGroup.GET("/public-key", api_whitelabel.WhitelabelGetPublicKey, middleware.Logging)
			whitelabelGroup.POST("/public-key", api_whitelabel.WhitelabelPostPublicKey, middleware.Logging)
			whitelabelGroup.POST("/create-interactions", api_whitelabel.GetWhitelabelCreateInteractions(), middleware.Logging)

			whitelabelGroup.POST("/", rl(middleware.RateLimitTypeUser, 10, time.Minute), api_whitelabel.WhitelabelPost, middleware.Logging)
			whitelabelGroup.POST("/status", rl(middleware.RateLimitTypeUser, 1, time.Second*5), api_whitelabel.WhitelabelStatusPost, middleware.Logging)
		}
	}

	if err := router.Run(config.Conf.Server.Host); err != nil {
		panic(err)
	}
}

func rl(rlType middleware.RateLimitType, limit int, period time.Duration) func(*gin.Context) {
	return middleware.CreateRateLimiter(rlType, limit, period)
}
