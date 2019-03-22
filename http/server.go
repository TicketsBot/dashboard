package http

import (
	"github.com/TicketsBot/PanelV2/config"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
)

func StartServer() {
	log.Println("Starting HTTP server")

	router := routing.New()

	err := fasthttp.ListenAndServe(config.Conf.Server.Host, router.HandleRequest); if err != nil {
		panic(err)
	}
}
