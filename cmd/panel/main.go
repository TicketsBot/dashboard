package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/manage"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/messagequeue"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/apex/log"
	gdlratelimit "github.com/rxdn/gdl/rest/ratelimit"
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
	database.ConnectToDatabase()
	cache.Instance = cache.NewCache()

	utils.LoadEmoji()

	messagequeue.Client = messagequeue.NewRedisClient()
	go Listen(messagequeue.Client)

	ratelimit.Ratelimiter = gdlratelimit.NewRateLimiter(gdlratelimit.NewRedisStore(messagequeue.Client.Client, "ratelimit")) // TODO: Use values from config

	http.StartServer()
}

func Listen(client messagequeue.RedisClient) {
	ch := make(chan messagequeue.TicketMessage)
	go client.ListenForMessages(ch)

	for decoded := range ch {
		manage.SocketsLock.Lock()
		for _, socket := range manage.Sockets {
			if socket.Guild == decoded.GuildId && socket.Ticket == decoded.TicketId {
				if err := socket.Ws.WriteJSON(decoded); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		manage.SocketsLock.Unlock()
	}
}
