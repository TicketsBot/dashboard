package http

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/api"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/manage"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/root"
	"github.com/TicketsBot/GoPanel/app/http/middleware"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/sessions"
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
	store, err := sessions.NewRedisStore(
		config.Conf.Server.Session.Threads,
		"tcp", fmt.Sprintf("%s:%d", config.Conf.Redis.Host, config.Conf.Redis.Port),
		config.Conf.Redis.Password,
		[]byte(config.Conf.Server.Session.Secret))
	if err != nil {
		panic(err)
	}
	router.Use(sessions.Sessions("panel", store))

	// Handle static asset requests
	router.Use(static.Serve("/assets/", static.LocalFile("./public/static", false)))

	router.Use(gin.Recovery())
	router.Use(createLimiter(600, time.Minute * 10))

	// Register templates
	router.HTMLRender = createRenderer()

	router.GET("/login", root.LoginHandler)
	router.GET("/callback", root.CallbackHandler)

	router.GET("/manage/:id/logs/view/:ticket", manage.LogViewHandler)              // we check in the actual handler bc of a custom redirect
	router.GET("/manage/:id/logs/modmail/view/:uuid", manage.ModmailLogViewHandler) // we check in the actual handler bc of a custom redirect

	authorized := router.Group("/", middleware.AuthenticateCookie)
	{
		authorized.POST("/token", api.TokenHandler)

		authenticateGuild := authorized.Group("/", middleware.AuthenticateGuild(false))

		authorized.GET("/", root.IndexHandler)
		authorized.GET("/whitelabel", root.WhitelabelHandler)
		authorized.GET("/logout", root.LogoutHandler)

		authenticateGuild.GET("/manage/:id/settings", manage.SettingsHandler)
		authenticateGuild.GET("/manage/:id/logs", manage.LogsHandler)
		authenticateGuild.GET("/manage/:id/logs/modmail", manage.ModmailLogsHandler)
		authenticateGuild.GET("/manage/:id/blacklist", manage.BlacklistHandler)
		authenticateGuild.GET("/manage/:id/panels", manage.PanelHandler)
		authenticateGuild.GET("/manage/:id/tags", manage.TagsHandler)

		authenticateGuild.GET("/manage/:id/tickets", manage.TicketListHandler)
		authenticateGuild.GET("/manage/:id/tickets/view/:ticketId", manage.TicketViewHandler)

		authorized.GET("/webchat", manage.WebChatWs)
	}

	apiGroup := router.Group("/api", middleware.AuthenticateToken)
	guildAuthApi := apiGroup.Group("/:id", middleware.AuthenticateGuild(true))
	{
		guildAuthApi.GET("/channels", api.ChannelsHandler)
		guildAuthApi.GET("/premium", api.PremiumHandler)
		guildAuthApi.GET("/user/:user", api.UserHandler)

		guildAuthApi.GET("/settings", api.GetSettingsHandler)
		guildAuthApi.POST("/settings", api.UpdateSettingsHandler)

		guildAuthApi.GET("/blacklist", api.GetBlacklistHandler)
		guildAuthApi.PUT("/blacklist", api.AddBlacklistHandler)
		guildAuthApi.DELETE("/blacklist/:user", api.RemoveBlacklistHandler)

		guildAuthApi.GET("/panels", api.ListPanels)
		guildAuthApi.PUT("/panels", api.CreatePanel)
		guildAuthApi.DELETE("/panels/:message", api.DeletePanel)

		guildAuthApi.GET("/logs/", api.GetLogs)
		guildAuthApi.GET("/modmail/logs/", api.GetModmailLogs)

		guildAuthApi.GET("/tickets", api.GetTickets)
		guildAuthApi.GET("/tickets/:ticketId", api.GetTicket)
		guildAuthApi.POST("/tickets/:ticketId", api.SendMessage)
		guildAuthApi.DELETE("/tickets/:ticketId", api.CloseTicket)

		guildAuthApi.GET("/tags", api.TagsListHandler)
		guildAuthApi.PUT("/tags", api.CreateTag)
		guildAuthApi.DELETE("/tags/:tag", api.DeleteTag)

		guildAuthApi.GET("/claimsettings", api.GetClaimSettings)
		guildAuthApi.POST("/claimsettings", api.PostClaimSettings)

		guildAuthApi.GET("/autoclose", api.GetAutoClose)
		guildAuthApi.POST("/autoclose", api.PostAutoClose)
	}

	userGroup := router.Group("/user", middleware.AuthenticateToken)
	{
		userGroup.GET("/guilds", api.GetGuilds)

		userGroup.GET("/whitelabel", api.WhitelabelGet)

		userGroup.Group("/").Use(createLimiter(10, time.Minute)).POST("/whitelabel", api.WhitelabelPost)
		userGroup.Group("/").Use(createLimiter(1, time.Second * 5)).POST("/whitelabel/status", api.WhitelabelStatusPost)
	}

	if err := router.Run(config.Conf.Server.Host); err != nil {
		panic(err)
	}
}

func createRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	r = addMainTemplate(r, "index")
	r = addMainTemplate(r, "whitelabel")

	r = addManageTemplate(r, "blacklist")
	r = addManageTemplate(r, "logs")
	r = addManageTemplate(r, "modmaillogs")
	r = addManageTemplate(r, "settings")
	r = addManageTemplate(r, "ticketlist")
	r = addManageTemplate(r, "ticketview")
	r = addManageTemplate(r, "panels")
	r = addManageTemplate(r, "tags")

	return r
}

func addMainTemplate(renderer multitemplate.Renderer, name string) multitemplate.Renderer {
	renderer.AddFromFiles(fmt.Sprintf("main/%s", name),
		"./public/templates/layouts/main.tmpl",
		"./public/templates/includes/head.tmpl",
		"./public/templates/includes/sidebar.tmpl",
		fmt.Sprintf("./public/templates/views/%s.tmpl", name),
	)
	return renderer
}

func addManageTemplate(renderer multitemplate.Renderer, name string) multitemplate.Renderer {
	renderer.AddFromFiles(fmt.Sprintf("manage/%s", name),
		"./public/templates/layouts/manage.tmpl",
		"./public/templates/includes/head.tmpl",
		"./public/templates/includes/sidebar.tmpl",
		"./public/templates/includes/navbar.tmpl",
		fmt.Sprintf("./public/templates/views/%s.tmpl", name),
	)
	return renderer
}

func createLimiter(limit int64, period time.Duration) func(*gin.Context) {
	store := memory.NewStore()
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}

	return mgin.NewMiddleware(limiter.New(store, rate))
}
