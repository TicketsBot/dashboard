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
	router.Use(createLimiter())

	// Register templates
	router.HTMLRender = createRenderer()

	router.GET("/login", root.LoginHandler)
	router.GET("/callback", root.CallbackHandler)

	router.GET("/manage/:id/logs/view/:ticket", manage.LogViewHandler) // we check in the actual handler bc of a custom redirect
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
	guildAuthApi := apiGroup.Group("/", middleware.AuthenticateGuild(true))
	{
		guildAuthApi.GET("/:id/channels", api.ChannelsHandler)
		guildAuthApi.GET("/:id/premium", api.PremiumHandler)
		guildAuthApi.GET("/:id/user/:user", api.UserHandler)

		guildAuthApi.GET("/:id/settings", api.GetSettingsHandler)
		guildAuthApi.POST("/:id/settings", api.UpdateSettingsHandler)

		guildAuthApi.GET("/:id/blacklist", api.GetBlacklistHandler)
		guildAuthApi.PUT("/:id/blacklist", api.AddBlacklistHandler)
		guildAuthApi.DELETE("/:id/blacklist/:user", api.RemoveBlacklistHandler)

		guildAuthApi.GET("/:id/panels", api.ListPanels)
		guildAuthApi.PUT("/:id/panels", api.CreatePanel)
		guildAuthApi.DELETE("/:id/panels/:message", api.DeletePanel)

		guildAuthApi.GET("/:id/logs/", api.GetLogs)
		guildAuthApi.GET("/:id/modmail/logs/", api.GetModmailLogs)

		guildAuthApi.GET("/:id/tickets", api.GetTickets)
		guildAuthApi.GET("/:id/tickets/:ticketId", api.GetTicket)
		guildAuthApi.POST("/:id/tickets/:ticketId", api.SendMessage)
		guildAuthApi.DELETE("/:id/tickets/:ticketId", api.CloseTicket)

		guildAuthApi.GET("/:id/tags", api.TagsListHandler)
		guildAuthApi.PUT("/:id/tags", api.CreateTag)
		guildAuthApi.DELETE("/:id/tags/:tag", api.DeleteTag)

		guildAuthApi.GET("/:id/claimsettings", api.GetClaimSettings)
		guildAuthApi.POST("/:id/claimsettings", api.PostClaimSettings)
	}

	userGroup := router.Group("/user", middleware.AuthenticateToken)
	{
		userGroup.GET("/guilds", api.GetGuilds)
		userGroup.POST("/whitelabel", api.WhitelabelHandler)
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

func createLimiter() func(*gin.Context) {
	store := memory.NewStore()
	rate := limiter.Rate{
		Period:    time.Minute * 10,
		Limit:     600,
	}

	return mgin.NewMiddleware(limiter.New(store, rate))
}
