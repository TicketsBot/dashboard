package http

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/api"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/manage"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/root"
	"github.com/TicketsBot/GoPanel/app/http/middleware"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/common/permission"
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

		authenticateGuildAdmin := authorized.Group("/", middleware.AuthenticateGuild(false, permission.Admin))
		authenticateGuildSupport := authorized.Group("/", middleware.AuthenticateGuild(false, permission.Support))

		authorized.GET("/", root.IndexHandler)
		authorized.GET("/whitelabel", root.WhitelabelHandler)
		authorized.GET("/logout", root.LogoutHandler)

		authenticateGuildAdmin.GET("/manage/:id/settings", manage.SettingsHandler)
		authenticateGuildSupport.GET("/manage/:id/logs", manage.LogsHandler)
		authenticateGuildSupport.GET("/manage/:id/logs/modmail", manage.ModmailLogsHandler)
		authenticateGuildSupport.GET("/manage/:id/blacklist", manage.BlacklistHandler)
		authenticateGuildAdmin.GET("/manage/:id/panels", manage.PanelHandler)
		authenticateGuildSupport.GET("/manage/:id/tags", manage.TagsHandler)

		authenticateGuildSupport.GET("/manage/:id/tickets", manage.TicketListHandler)
		authenticateGuildSupport.GET("/manage/:id/tickets/view/:ticketId", manage.TicketViewHandler)

		authorized.GET("/webchat", manage.WebChatWs)
	}

	apiGroup := router.Group("/api", middleware.AuthenticateToken)
	guildAuthApiAdmin := apiGroup.Group("/:id", middleware.AuthenticateGuild(true, permission.Admin))
	guildAuthApiSupport := apiGroup.Group("/:id", middleware.AuthenticateGuild(true, permission.Support))
	{
		guildAuthApiSupport.GET("/channels", api.ChannelsHandler)
		guildAuthApiSupport.GET("/premium", api.PremiumHandler)
		guildAuthApiSupport.GET("/user/:user", api.UserHandler)
		guildAuthApiSupport.GET("/roles", api.RolesHandler)

		guildAuthApiAdmin.GET("/settings", api.GetSettingsHandler)
		guildAuthApiAdmin.POST("/settings", api.UpdateSettingsHandler)

		guildAuthApiSupport.GET("/blacklist", api.GetBlacklistHandler)
		guildAuthApiSupport.PUT("/blacklist", api.AddBlacklistHandler)
		guildAuthApiSupport.DELETE("/blacklist/:user", api.RemoveBlacklistHandler)

		guildAuthApiAdmin.GET("/panels", api.ListPanels)
		guildAuthApiAdmin.PUT("/panels", api.CreatePanel)
		guildAuthApiAdmin.PUT("/panels/:message", api.UpdatePanel)
		guildAuthApiAdmin.DELETE("/panels/:message", api.DeletePanel)

		guildAuthApiSupport.GET("/logs/", api.GetLogs)
		guildAuthApiSupport.GET("/modmail/logs/", api.GetModmailLogs)

		guildAuthApiSupport.GET("/tickets", api.GetTickets)
		guildAuthApiSupport.GET("/tickets/:ticketId", api.GetTicket)
		guildAuthApiSupport.POST("/tickets/:ticketId", api.SendMessage)
		guildAuthApiSupport.DELETE("/tickets/:ticketId", api.CloseTicket)

		guildAuthApiSupport.GET("/tags", api.TagsListHandler)
		guildAuthApiSupport.PUT("/tags", api.CreateTag)
		guildAuthApiSupport.DELETE("/tags/:tag", api.DeleteTag)

		guildAuthApiAdmin.GET("/claimsettings", api.GetClaimSettings)
		guildAuthApiAdmin.POST("/claimsettings", api.PostClaimSettings)

		guildAuthApiAdmin.GET("/autoclose", api.GetAutoClose)
		guildAuthApiAdmin.POST("/autoclose", api.PostAutoClose)
	}

	userGroup := router.Group("/user", middleware.AuthenticateToken)
	{
		userGroup.GET("/guilds", api.GetGuilds)
		userGroup.GET("/permissionlevel", api.GetPermissionLevel)

		{
			whitelabelGroup := userGroup.Group("/whitelabel", middleware.VerifyWhitelabel(false))
			whitelabelApiGroup := userGroup.Group("/whitelabel", middleware.VerifyWhitelabel(true))

			whitelabelGroup.GET("/", api.WhitelabelGet)
			whitelabelApiGroup.GET("/errors", api.WhitelabelGetErrors)
			whitelabelApiGroup.GET("/guilds", api.WhitelabelGetGuilds)
			whitelabelApiGroup.POST("/modmail", api.WhitelabelModmailPost)

			whitelabelApiGroup.Group("/").Use(createLimiter(10, time.Minute)).POST("/", api.WhitelabelPost)
			whitelabelApiGroup.Group("/").Use(createLimiter(1, time.Second * 5)).POST("/status", api.WhitelabelStatusPost)
		}
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
	r = addManageTemplate(r, "settings", "./public/templates/includes/substitutionmodal.tmpl")
	r = addManageTemplate(r, "ticketlist")
	r = addManageTemplate(r, "ticketview")
	r = addManageTemplate(r, "panels", "./public/templates/includes/substitutionmodal.tmpl", "./public/templates/includes/paneleditmodal.tmpl")
	r = addManageTemplate(r, "tags")

	r = addErrorTemplate(r)

	return r
}

func addMainTemplate(renderer multitemplate.Renderer, name string, extra ...string) multitemplate.Renderer {
	files := []string{
		"./public/templates/layouts/main.tmpl",
		"./public/templates/includes/head.tmpl",
		"./public/templates/includes/sidebar.tmpl",
		"./public/templates/includes/loadingscreen.tmpl",
		fmt.Sprintf("./public/templates/views/%s.tmpl", name),
	}

	files = append(files, extra...)

	renderer.AddFromFiles(fmt.Sprintf("main/%s", name), files...)
	return renderer
}

func addManageTemplate(renderer multitemplate.Renderer, name string, extra ...string) multitemplate.Renderer {
	files := []string{
		"./public/templates/layouts/manage.tmpl",
		"./public/templates/includes/head.tmpl",
		"./public/templates/includes/sidebar.tmpl",
		"./public/templates/includes/navbar.tmpl",
		"./public/templates/includes/loadingscreen.tmpl",
		fmt.Sprintf("./public/templates/views/%s.tmpl", name),
	}

	files = append(files, extra...)

	renderer.AddFromFiles(fmt.Sprintf("manage/%s", name), files...)
	return renderer
}

func addErrorTemplate(renderer multitemplate.Renderer) multitemplate.Renderer {
	files := []string{
		"./public/templates/layouts/error.tmpl",
		"./public/templates/includes/head.tmpl",
	}

	renderer.AddFromFiles("error", files...)
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
