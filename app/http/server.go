package http

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/api"
	api_autoclose "github.com/TicketsBot/GoPanel/app/http/endpoints/api/autoclose"
	api_blacklist "github.com/TicketsBot/GoPanel/app/http/endpoints/api/blacklist"
	api_logs "github.com/TicketsBot/GoPanel/app/http/endpoints/api/logs"
	api_panels "github.com/TicketsBot/GoPanel/app/http/endpoints/api/panel"
	api_settings "github.com/TicketsBot/GoPanel/app/http/endpoints/api/settings"
	api_tags "github.com/TicketsBot/GoPanel/app/http/endpoints/api/tags"
	api_team "github.com/TicketsBot/GoPanel/app/http/endpoints/api/team"
	api_ticket "github.com/TicketsBot/GoPanel/app/http/endpoints/api/ticket"
	api_whitelabel "github.com/TicketsBot/GoPanel/app/http/endpoints/api/whitelabel"
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
	router.Use(createLimiter(600, time.Minute*10))

	// Register templates
	router.HTMLRender = createRenderer()

	router.GET("/login", root.LoginHandler)
	router.GET("/callback", root.CallbackHandler)

	router.GET("/manage/:id/logs/view/:ticket", manage.LogViewHandler) // we check in the actual handler bc of a custom redirect

	authorized := router.Group("/", middleware.AuthenticateCookie)
	{
		authorized.POST("/token", middleware.VerifyXTicketsHeader, api.TokenHandler)

		authenticateGuildAdmin := authorized.Group("/", middleware.AuthenticateGuild(false, permission.Admin))
		authenticateGuildSupport := authorized.Group("/", middleware.AuthenticateGuild(false, permission.Support))

		authorized.GET("/", root.IndexHandler)
		authorized.GET("/whitelabel", root.WhitelabelHandler)
		authorized.GET("/logout", root.LogoutHandler)

		authenticateGuildAdmin.GET("/manage/:id/settings", manage.SettingsHandler)
		authenticateGuildSupport.GET("/manage/:id/logs", manage.LogsHandler)
		authenticateGuildSupport.GET("/manage/:id/blacklist", manage.BlacklistHandler)
		authenticateGuildAdmin.GET("/manage/:id/panels", manage.PanelHandler)
		authenticateGuildSupport.GET("/manage/:id/tags", manage.TagsHandler)
		authenticateGuildSupport.GET("/manage/:id/teams", serveTemplate("manage/teams"))

		authenticateGuildSupport.GET("/manage/:id/tickets", manage.TicketListHandler)
		authenticateGuildSupport.GET("/manage/:id/tickets/view/:ticketId", manage.TicketViewHandler)

		authorized.GET("/webchat", manage.WebChatWs)
	}

	apiGroup := router.Group("/api", middleware.VerifyXTicketsHeader, middleware.AuthenticateToken)
	guildAuthApiAdmin := apiGroup.Group("/:id", middleware.AuthenticateGuild(true, permission.Admin))
	guildAuthApiSupport := apiGroup.Group("/:id", middleware.AuthenticateGuild(true, permission.Support))
	{
		guildAuthApiSupport.GET("/channels", api.ChannelsHandler)
		guildAuthApiSupport.GET("/premium", api.PremiumHandler)
		guildAuthApiSupport.GET("/user/:user", api.UserHandler)
		guildAuthApiSupport.GET("/roles", api.RolesHandler)
		guildAuthApiSupport.GET("/members/search", createLimiter(10, time.Second * 30), createLimiter(75, time.Minute * 30), api.SearchMembers)

		guildAuthApiAdmin.GET("/settings", api_settings.GetSettingsHandler)
		guildAuthApiAdmin.POST("/settings", api_settings.UpdateSettingsHandler)

		guildAuthApiSupport.GET("/blacklist", api_blacklist.GetBlacklistHandler)
		guildAuthApiSupport.PUT("/blacklist", api_blacklist.AddBlacklistHandler)
		guildAuthApiSupport.DELETE("/blacklist/:user", api_blacklist.RemoveBlacklistHandler)

		guildAuthApiAdmin.GET("/panels", api_panels.ListPanels)
		guildAuthApiAdmin.PUT("/panels", api_panels.CreatePanel)
		guildAuthApiAdmin.PUT("/panels/:message", api_panels.UpdatePanel)
		guildAuthApiAdmin.DELETE("/panels/:message", api_panels.DeletePanel)

		guildAuthApiAdmin.GET("/multipanels", api_panels.MultiPanelList)
		guildAuthApiAdmin.POST("/multipanels", api_panels.MultiPanelCreate)
		guildAuthApiAdmin.PATCH("/multipanels/:panelid", api_panels.MultiPanelUpdate)
		guildAuthApiAdmin.DELETE("/multipanels/:panelid", api_panels.MultiPanelDelete)

		guildAuthApiSupport.GET("/logs/", api_logs.GetLogs)

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
		guildAuthApiAdmin.GET("/team/:teamid", createLimiter(5, time.Second * 15), api_team.GetMembers)
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
			whitelabelGroup := userGroup.Group("/whitelabel", middleware.VerifyWhitelabel(false))
			whitelabelApiGroup := userGroup.Group("/whitelabel", middleware.VerifyWhitelabel(true))

			whitelabelGroup.GET("/", api_whitelabel.WhitelabelGet)
			whitelabelApiGroup.GET("/errors", api_whitelabel.WhitelabelGetErrors)
			whitelabelApiGroup.GET("/guilds", api_whitelabel.WhitelabelGetGuilds)
			whitelabelApiGroup.GET("/public-key", api_whitelabel.WhitelabelGetPublicKey)
			whitelabelApiGroup.POST("/public-key", api_whitelabel.WhitelabelPostPublicKey)
			whitelabelApiGroup.POST("/create-interactions", api_whitelabel.WhitelabelCreateInteractions)

			whitelabelApiGroup.POST("/", createLimiter(10, time.Minute), api_whitelabel.WhitelabelPost)
			whitelabelApiGroup.POST("/status", createLimiter(1, time.Second*5), api_whitelabel.WhitelabelStatusPost)
		}
	}

	if err := router.Run(config.Conf.Server.Host); err != nil {
		panic(err)
	}
}

func serveTemplate(templateName string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		store := sessions.Default(ctx)
		guildId := ctx.Keys["guildid"].(uint64)

		ctx.HTML(200, templateName, gin.H{
			"name":         store.Get("name").(string),
			"guildId":      guildId,
			"avatar":       store.Get("avatar").(string),
			"baseUrl":      config.Conf.Server.BaseUrl,
		})
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
	r = addManageTemplate(r, "panels", "./public/templates/includes/substitutionmodal.tmpl", "./public/templates/includes/paneleditmodal.tmpl", "./public/templates/includes/multipaneleditmodal.tmpl")
	r = addManageTemplate(r, "tags")
	r = addManageTemplate(r, "teams")

	r = addErrorTemplate(r)

	return r
}

func addMainTemplate(renderer multitemplate.Renderer, name string, extra ...string) multitemplate.Renderer {
	files := []string{
		"./public/templates/layouts/main.tmpl",
		"./public/templates/includes/head.tmpl",
		"./public/templates/includes/sidebar.tmpl",
		"./public/templates/includes/loadingscreen.tmpl",
		"./public/templates/includes/notifymodal.tmpl",
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
		"./public/templates/includes/notifymodal.tmpl",
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
