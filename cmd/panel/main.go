package main

import (
	"github.com/TicketsBot/GoPanel/app/http"
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
	go cache.Client.ListenForMessages()

	http.StartServer()
}
