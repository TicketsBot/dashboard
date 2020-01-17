package main

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/manage"
	"github.com/TicketsBot/GoPanel/cache"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano() % 3497)

	config.LoadConfig()
	database.ConnectToDatabase()

	cache.Client = cache.NewRedisClient()
	go Listen(cache.Client)

	http.StartServer()
}

func Listen(client cache.RedisClient) {
	ch := make(chan cache.TicketMessage)
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
