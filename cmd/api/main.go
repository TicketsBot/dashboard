package main

import (
	"fmt"
	app "github.com/TicketsBot/GoPanel/app/http"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/api/ticket/livechat"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/archiverclient"
	"github.com/TicketsBot/common/chatrelay"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/common/secureproxy"
	"github.com/TicketsBot/worker/i18n"
	"github.com/getsentry/sentry-go"
	"github.com/rxdn/gdl/rest/request"
	"net/http"
	"net/http/pprof"
	"time"
)

func main() {
	startPprof()

	config.LoadConfig()

	sentryOpts := sentry.ClientOptions{
		Dsn:              config.Conf.SentryDsn,
		Debug:            config.Conf.Debug,
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: 0.1,
	}

	if err := sentry.Init(sentryOpts); err != nil {
		fmt.Printf("Error initialising sentry: %s", err.Error())
	}

	fmt.Println("Connecting to database...")
	database.ConnectToDatabase()

	fmt.Println("Connecting to cache...")
	cache.Instance = cache.NewCache()

	fmt.Println("Initialising microservice clients...")
	utils.ArchiverClient = archiverclient.NewArchiverClientWithTimeout(config.Conf.Bot.ObjectStore, time.Second*15, []byte(config.Conf.Bot.AesKey))
	utils.SecureProxyClient = secureproxy.NewSecureProxy(config.Conf.SecureProxyUrl)

	utils.LoadEmoji()

	i18n.Init()

	if config.Conf.Bot.ProxyUrl != "" {
		request.RegisterHook(utils.ProxyHook)
	}

	fmt.Println("Connecting to Redis...")
	redis.Client = redis.NewRedisClient()

	socketManager := livechat.NewSocketManager()
	go socketManager.Run()

	go ListenChat(redis.Client, socketManager)

	if !config.Conf.Debug {
		rpc.PremiumClient = premium.NewPremiumLookupClient(
			premium.NewPatreonClient(config.Conf.Bot.PremiumLookupProxyUrl, config.Conf.Bot.PremiumLookupProxyKey),
			redis.Client.Client,
			cache.Instance.PgCache,
			database.Client,
		)
	} else {
		c := premium.NewMockLookupClient(premium.Whitelabel, premium.SourcePatreon)
		rpc.PremiumClient = &c
	}

	fmt.Println("Starting server...")
	app.StartServer(socketManager)
}

func ListenChat(client redis.RedisClient, sm *livechat.SocketManager) {
	ch := make(chan chatrelay.MessageData)
	go chatrelay.Listen(client.Client, ch)

	for event := range ch {
		sm.BroadcastMessage(event)
	}
}

func startPprof() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/{action}", pprof.Index)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	go func() {
		http.ListenAndServe(":6060", mux)
	}()
}
