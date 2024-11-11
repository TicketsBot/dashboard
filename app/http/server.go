package http

import (
	"github.com/TicketsBot/GoPanel/app/http/endpoints/api"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/api/admin/botstaff"
	api_blacklist "github.com/TicketsBot/GoPanel/app/http/endpoints/api/blacklist"
	api_forms "github.com/TicketsBot/GoPanel/app/http/endpoints/api/forms"
	api_integrations "github.com/TicketsBot/GoPanel/app/http/endpoints/api/integrations"
	api_panels "github.com/TicketsBot/GoPanel/app/http/endpoints/api/panel"
	api_premium "github.com/TicketsBot/GoPanel/app/http/endpoints/api/premium"
	api_settings "github.com/TicketsBot/GoPanel/app/http/endpoints/api/settings"
	api_override "github.com/TicketsBot/GoPanel/app/http/endpoints/api/staffoverride"
	api_tags "github.com/TicketsBot/GoPanel/app/http/endpoints/api/tags"
	api_team "github.com/TicketsBot/GoPanel/app/http/endpoints/api/team"
	api_ticket "github.com/TicketsBot/GoPanel/app/http/endpoints/api/ticket"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/api/ticket/livechat"
	api_transcripts "github.com/TicketsBot/GoPanel/app/http/endpoints/api/transcripts"
	api_whitelabel "github.com/TicketsBot/GoPanel/app/http/endpoints/api/whitelabel"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/root"
	"github.com/TicketsBot/GoPanel/app/http/middleware"
	"github.com/TicketsBot/GoPanel/app/http/session"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/common/permission"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.uber.org/zap"
	"time"
)

func StartServer(logger *zap.Logger, sm *livechat.SocketManager) {
	logger.Info("Starting HTTP server")

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logging(logger))

	router.RemoteIPHeaders = config.Conf.Server.RealIpHeaders
	if err := router.SetTrustedProxies(config.Conf.Server.TrustedProxies); err != nil {
		panic(err)
	}

	// Sessions
	session.Store = session.NewRedisStore()

	router.Use(rl(middleware.RateLimitTypeIp, 60, time.Minute))
	router.Use(rl(middleware.RateLimitTypeIp, 20, time.Second*10))
	router.Use(rl(middleware.RateLimitTypeUser, 60, time.Minute))
	router.Use(rl(middleware.RateLimitTypeGuild, 600, time.Minute*5))

	router.Use(middleware.Cors(config.Conf))

	// Metrics
	if len(config.Conf.Server.MetricHost) > 0 {
		monitor := ginmetrics.GetMonitor()
		monitor.UseWithoutExposingEndpoint(router)
		monitor.SetMetricPath("/metrics")

		metricRouter := gin.New()
		metricRouter.Use(gin.Recovery())
		metricRouter.Use(middleware.Logging(logger))
		
		monitor.Expose(metricRouter)

		go func() {
			if err := metricRouter.Run(config.Conf.Server.MetricHost); err != nil {
				panic(err)
			}
		}()
	}

	// util endpoints
	router.GET("/ip", root.IpHandler)
	router.GET("/robots.txt", func(ctx *gin.Context) {
		ctx.String(200, "Disallow: /")
	})

	router.POST("/callback", middleware.VerifyXTicketsHeader, root.CallbackHandler)
	router.POST("/logout", middleware.VerifyXTicketsHeader, middleware.AuthenticateToken, root.LogoutHandler)

	apiGroup := router.Group("/api", middleware.VerifyXTicketsHeader, middleware.AuthenticateToken, middleware.UpdateLastSeen)
	{
		apiGroup.GET("/session", api.SessionHandler)

		{
			integrationGroup := apiGroup.Group("/integrations")

			integrationGroup.GET("/self", api_integrations.GetOwnedIntegrationsHandler)
			integrationGroup.GET("/view/:integrationid", api_integrations.GetIntegrationHandler)
			integrationGroup.GET("/view/:integrationid/detail", api_integrations.GetIntegrationDetailedHandler)
			integrationGroup.POST("/:integrationid/public", api_integrations.SetIntegrationPublicHandler)
			integrationGroup.PATCH("/:integrationid", api_integrations.UpdateIntegrationHandler)
			integrationGroup.DELETE("/:integrationid", api_integrations.DeleteIntegrationHandler)
			apiGroup.POST("/integrations", api_integrations.CreateIntegrationHandler)
		}

		{
			premiumGroup := apiGroup.Group("/premium/@me")
			premiumGroup.GET("/entitlements", api_premium.GetEntitlements)
			premiumGroup.PUT("/active-guilds", api_premium.SetActiveGuilds)
		}
	}

	guildAuthApiAdmin := apiGroup.Group("/:id", middleware.AuthenticateGuild(permission.Admin))
	guildAuthApiSupport := apiGroup.Group("/:id", middleware.AuthenticateGuild(permission.Support))
	guildApiNoAuth := apiGroup.Group("/:id", middleware.ParseGuildId)
	{
		guildAuthApiSupport.GET("/guild", api.GuildHandler)
		guildAuthApiSupport.GET("/channels", api.ChannelsHandler)
		guildAuthApiSupport.GET("/premium", api.PremiumHandler)
		guildAuthApiSupport.GET("/user/:user", api.UserHandler)
		guildAuthApiSupport.GET("/roles", api.RolesHandler)
		guildAuthApiSupport.GET("/emojis", rl(middleware.RateLimitTypeGuild, 5, time.Second*30), api.EmojisHandler)
		guildAuthApiSupport.GET("/members/search",
			rl(middleware.RateLimitTypeGuild, 5, time.Second),
			rl(middleware.RateLimitTypeGuild, 10, time.Second*30),
			rl(middleware.RateLimitTypeGuild, 75, time.Minute*30),
			api.SearchMembers,
		)

		// Must be readable to load transcripts page
		guildAuthApiSupport.GET("/settings", api_settings.GetSettingsHandler)
		guildAuthApiAdmin.POST("/settings", api_settings.UpdateSettingsHandler)

		guildAuthApiSupport.GET("/blacklist", api_blacklist.GetBlacklistHandler)
		guildAuthApiSupport.POST("/blacklist", api_blacklist.AddBlacklistHandler)
		guildAuthApiSupport.DELETE("/blacklist/user/:user", api_blacklist.RemoveUserBlacklistHandler)
		guildAuthApiSupport.DELETE("/blacklist/role/:role", api_blacklist.RemoveRoleBlacklistHandler)

		// Must be readable to load transcripts page
		guildAuthApiSupport.GET("/panels", api_panels.ListPanels)
		guildAuthApiAdmin.POST("/panels", api_panels.CreatePanel)
		guildAuthApiAdmin.POST("/panels/:panelid", rl(middleware.RateLimitTypeGuild, 5, 5*time.Second), api_panels.ResendPanel)
		guildAuthApiAdmin.PATCH("/panels/:panelid", api_panels.UpdatePanel)
		guildAuthApiAdmin.DELETE("/panels/:panelid", api_panels.DeletePanel)

		guildAuthApiAdmin.GET("/multipanels", api_panels.MultiPanelList)
		guildAuthApiAdmin.POST("/multipanels", api_panels.MultiPanelCreate)
		guildAuthApiAdmin.POST("/multipanels/:panelid", rl(middleware.RateLimitTypeGuild, 5, 5*time.Second), api_panels.MultiPanelResend)
		guildAuthApiAdmin.PATCH("/multipanels/:panelid", api_panels.MultiPanelUpdate)
		guildAuthApiAdmin.DELETE("/multipanels/:panelid", api_panels.MultiPanelDelete)

		guildAuthApiSupport.GET("/forms", api_forms.GetForms)
		guildAuthApiAdmin.POST("/forms", rl(middleware.RateLimitTypeGuild, 30, time.Hour), api_forms.CreateForm)
		guildAuthApiAdmin.PATCH("/forms/:form_id", rl(middleware.RateLimitTypeGuild, 30, time.Hour), api_forms.UpdateForm)
		guildAuthApiAdmin.DELETE("/forms/:form_id", api_forms.DeleteForm)
		guildAuthApiAdmin.PATCH("/forms/:form_id/inputs", api_forms.UpdateInputs)

		// Should be a GET, but easier to take a body for development purposes
		guildAuthApiSupport.POST("/transcripts",
			rl(middleware.RateLimitTypeUser, 5, 5*time.Second),
			rl(middleware.RateLimitTypeUser, 20, time.Minute),
			api_transcripts.ListTranscripts,
		)

		// Allow regular users to get their own transcripts, make sure you check perms inside
		guildApiNoAuth.GET("/transcripts/:ticketId", rl(middleware.RateLimitTypeGuild, 10, 10*time.Second), api_transcripts.GetTranscriptHandler)
		guildApiNoAuth.GET("/transcripts/:ticketId/render", rl(middleware.RateLimitTypeGuild, 10, 10*time.Second), api_transcripts.GetTranscriptRenderHandler)

		guildAuthApiSupport.GET("/tickets", api_ticket.GetTickets)
		guildAuthApiSupport.GET("/tickets/:ticketId", api_ticket.GetTicket)
		guildAuthApiSupport.POST("/tickets/:ticketId", rl(middleware.RateLimitTypeGuild, 5, time.Second*5), api_ticket.SendMessage)
		guildAuthApiSupport.POST("/tickets/:ticketId/tag", rl(middleware.RateLimitTypeGuild, 5, time.Second*5), api_ticket.SendTag)
		guildAuthApiSupport.DELETE("/tickets/:ticketId", api_ticket.CloseTicket)

		// Websockets do not support headers: so we must implement authentication over the WS connection
		router.GET("/api/:id/tickets/:ticketId/live-chat", livechat.GetLiveChatHandler(sm))

		guildAuthApiSupport.GET("/tags", api_tags.TagsListHandler)
		guildAuthApiSupport.PUT("/tags", api_tags.CreateTag)
		guildAuthApiSupport.DELETE("/tags", api_tags.DeleteTag)

		guildAuthApiAdmin.GET("/team", api_team.GetTeams)
		guildAuthApiAdmin.GET("/team/:teamid", rl(middleware.RateLimitTypeUser, 10, time.Second*30), api_team.GetMembers)
		guildAuthApiAdmin.POST("/team", rl(middleware.RateLimitTypeUser, 10, time.Minute), api_team.CreateTeam)
		guildAuthApiAdmin.PUT("/team/:teamid/:snowflake", rl(middleware.RateLimitTypeGuild, 5, time.Second*10), api_team.AddMember)
		guildAuthApiAdmin.DELETE("/team/:teamid", api_team.DeleteTeam)
		guildAuthApiAdmin.DELETE("/team/:teamid/:snowflake", rl(middleware.RateLimitTypeGuild, 30, time.Minute), api_team.RemoveMember)

		guildAuthApiAdmin.GET("/staff-override", api_override.GetOverrideHandler)
		guildAuthApiAdmin.POST("/staff-override", api_override.CreateOverrideHandler)
		guildAuthApiAdmin.DELETE("/staff-override", api_override.DeleteOverrideHandler)

		guildAuthApiAdmin.GET("/integrations/available", api_integrations.ListIntegrationsHandler)
		guildAuthApiAdmin.GET("/integrations/:integrationid", api_integrations.IsIntegrationActiveHandler)
		guildAuthApiAdmin.POST("/integrations/:integrationid",
			rl(middleware.RateLimitTypeUser, 10, time.Minute),
			rl(middleware.RateLimitTypeGuild, 10, time.Minute),
			rl(middleware.RateLimitTypeUser, 30, time.Minute*30),
			rl(middleware.RateLimitTypeGuild, 30, time.Minute*30),
			api_integrations.ActivateIntegrationHandler,
		)
		guildAuthApiAdmin.PATCH("/integrations/:integrationid", api_integrations.UpdateIntegrationSecretsHandler)
		guildAuthApiAdmin.DELETE("/integrations/:integrationid", api_integrations.RemoveIntegrationHandler)
	}

	userGroup := router.Group("/user", middleware.AuthenticateToken, middleware.UpdateLastSeen)
	{
		userGroup.GET("/guilds", api.GetGuilds)
		userGroup.POST("/guilds/reload", api.ReloadGuildsHandler)
		userGroup.GET("/permissionlevel", api.GetPermissionLevel)

		{
			whitelabelGroup := userGroup.Group("/whitelabel", middleware.VerifyWhitelabel(true))

			whitelabelGroup.GET("/", api_whitelabel.WhitelabelGet)
			whitelabelGroup.GET("/errors", api_whitelabel.WhitelabelGetErrors)
			whitelabelGroup.GET("/guilds", api_whitelabel.WhitelabelGetGuilds)
			whitelabelGroup.POST("/create-interactions", api_whitelabel.GetWhitelabelCreateInteractions())
			whitelabelGroup.DELETE("/", api_whitelabel.WhitelabelDelete)

			whitelabelGroup.POST("/", rl(middleware.RateLimitTypeUser, 5, time.Minute), api_whitelabel.WhitelabelPost())
			whitelabelGroup.POST("/status", rl(middleware.RateLimitTypeUser, 1, time.Second*5), api_whitelabel.WhitelabelStatusPost)
		}
	}

	adminGroup := apiGroup.Group("/admin", middleware.AdminOnly)
	{
		adminGroup.GET("/bot-staff", botstaff.ListBotStaffHandler)
		adminGroup.POST("/bot-staff/:userid", botstaff.AddBotStaffHandler)
		adminGroup.DELETE("/bot-staff/:userid", botstaff.RemoveBotStaffHandler)
	}

	if err := router.Run(config.Conf.Server.Host); err != nil {
		panic(err)
	}
}

func rl(rlType middleware.RateLimitType, limit int, period time.Duration) func(*gin.Context) {
	return middleware.CreateRateLimiter(rlType, limit, period)
}
