package http

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/manage"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/root"
	"github.com/TicketsBot/GoPanel/app/http/template"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

func StartServer() {
	log.Println("Starting HTTP server")

	// Compile templates
	template.LoadLayouts()
	template.LoadTemplates()

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

	// Root
	router.GET("/", root.IndexHandler)

	// /login
	router.GET("/login", root.LoginHandler)

	// /callback
	router.GET("/callback", root.CallbackHandler)

	// /logout
	router.GET("/logout", root.LogoutHandler)

	// /manage/:id/settings
	router.GET("/manage/:id/settings", manage.SettingsHandler)

	// /manage/:id/logs/page/:page
	router.GET("/manage/:id/logs/page/:page", manage.LogsHandler)

	// /manage/:id/logs/view/:uuid
	router.GET("/manage/:id/logs/view/:uuid", manage.LogViewHandler)

	// /manage/:id/blacklist
	router.GET("/manage/:id/blacklist", manage.BlacklistHandler)

	// /manage/:id/blacklist/remove/:user
	router.GET("/manage/:id/blacklist/remove/:user", manage.BlacklistRemoveHandler)

	if err := router.Run(config.Conf.Server.Host); err != nil {
		panic(err)
	}
}
