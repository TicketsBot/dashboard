package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/TicketsBot/GoPanel/app/http"
	"github.com/TicketsBot/GoPanel/app/http/endpoints/manage"
	"github.com/TicketsBot/GoPanel/cache"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/apex/log"
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

	utils.LoadEmoji()

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
