package http

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/manage"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/root"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
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

	// Register templates
	router.HTMLRender = createRenderer()

	router.GET("/", root.IndexHandler)

	router.GET("/login", root.LoginHandler)
	router.GET("/callback", root.CallbackHandler)
	router.GET("/logout", root.LogoutHandler)

	router.GET("/manage/:id/settings", manage.SettingsHandler)
	router.POST("/manage/:id/settings", manage.UpdateSettingsHandler)

	router.GET("/manage/:id/logs/page/:page", manage.LogsHandler)
	router.GET("/manage/:id/logs/view/:uuid", manage.LogViewHandler)

	router.GET("/manage/:id/blacklist", manage.BlacklistHandler)
	router.GET("/manage/:id/blacklist/remove/:user", manage.BlacklistRemoveHandler)

	router.GET("/manage/:id/panels", manage.PanelHandler)
	router.POST("/manage/:id/panels/create", manage.PanelCreateHandler)
	router.GET("/manage/:id/panels/delete/:msg", manage.PanelDeleteHandler)

	router.GET("/manage/:id/tickets", manage.TicketListHandler)
	router.GET("/manage/:id/tickets/view/:uuid", manage.TicketViewHandler)
	router.POST("/manage/:id/tickets/view/:uuid", manage.SendMessage)
	router.GET("/webchat", manage.WebChatWs)

	if err := router.Run(config.Conf.Server.Host); err != nil {
		panic(err)
	}
}

func createRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	r = addMainTemplate(r, "index")

	r = addManageTemplate(r, "blacklist")
	r = addManageTemplate(r, "logs")
	r = addManageTemplate(r, "settings")
	r = addManageTemplate(r, "ticketlist")
	r = addManageTemplate(r, "ticketview")
	r = addManageTemplate(r, "panels")

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

