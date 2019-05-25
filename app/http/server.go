package http

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http/endpoints"
	"github.com/TicketsBot/GoPanel/config"
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

	// Root
	router.GET("/", func(c *gin.Context) {
		endpoints.IndexHandler(c)
	})

	// /login
	router.GET("/login", func(c *gin.Context) {
		endpoints.LoginHandler(c)
	})

	// /callback
	router.GET("/callback", func(c *gin.Context) {
		endpoints.CallbackHandler(c)
	})

	if err := router.Run(config.Conf.Server.Host); err != nil {
		panic(err)
	}
}
