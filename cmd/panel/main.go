package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/root"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/redis"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/archiverclient"
	"github.com/TicketsBot/common/chatrelay"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/worker/i18n"
	"github.com/apex/log"
	"github.com/getsentry/sentry-go"
	"github.com/rxdn/gdl/rest/request"
	"math/rand"
	"time"
)

func main() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err == nil {
		rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	} else {
		log.Error(err.Error())
		rand.Seed(time.Now().UnixNano())
	}

	config.LoadConfig()

	sentryOpts := sentry.ClientOptions{
		Dsn:              config.Conf.SentryDsn,
		Debug:            config.Conf.Debug,
		AttachStacktrace: true,

	}
	if err := sentry.Init(sentryOpts); err != nil {
		fmt.Printf("Error initialising sentry: %s", err.Error())
	}

	database.ConnectToDatabase()

	cache.Instance = cache.NewCache()

	utils.ArchiverClient = archiverclient.NewArchiverClientWithTimeout(config.Conf.Bot.ObjectStore, time.Second*15, []byte(config.Conf.Bot.AesKey))

	utils.LoadEmoji()

	i18n.LoadMessages()

	if config.Conf.Bot.ProxyUrl != "" {
		request.RegisterHook(utils.ProxyHook)
	}

	redis.Client = redis.NewRedisClient()
	go ListenChat(redis.Client)

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

	http.StartServer()
}

func ListenChat(client redis.RedisClient) {
	ch := make(chan chatrelay.MessageData)
	go chatrelay.Listen(client.Client, ch)

	for event := range ch {
		root.SocketsLock.RLock()
		for _, socket := range root.Sockets {
			if socket.GuildId == event.Ticket.GuildId && socket.TicketId == event.Ticket.Id {
				if err := socket.Ws.WriteJSON(event.Message); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		root.SocketsLock.RUnlock()
	}
}
