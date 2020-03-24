package manage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord"
	"github.com/TicketsBot/GoPanel/utils/discord/endpoints/channel"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var SocketsLock sync.Mutex
var Sockets []*Socket

type (
	Socket struct {
		Ws     *websocket.Conn
		Guild  string
		Ticket int
	}

	WsEvent struct {
		Type string
		Data interface{}
	}

	AuthEvent struct {
		Guild  string
		Ticket string
	}
)

func WebChatWs(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	if utils.IsLoggedIn(store) {
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		socket := &Socket{
			Ws: conn,
		}

		conn.SetCloseHandler(func(code int, text string) error {
			i := -1
			SocketsLock.Lock()

			for index, element := range Sockets {
				if element == socket {
					i = index
					break
				}
			}

			if i != -1 {
				Sockets = Sockets[:i+copy(Sockets[i:], Sockets[i+1:])]
			}
			SocketsLock.Unlock()

			return nil
		})

		SocketsLock.Lock()
		Sockets = append(Sockets, socket)
		SocketsLock.Unlock()

		userIdStr := store.Get("userid").(string)
		userId, err := utils.GetUserId(store)
		if err != nil {
			conn.Close()
			return
		}

		var guildId string
		var guildIdParsed int64
		var ticket int

		for {
			var evnt WsEvent
			err := conn.ReadJSON(&evnt)
			if err != nil {
				break
			}

			if guildId == "" && evnt.Type != "auth" {
				conn.Close()
				break
			} else if evnt.Type == "auth" {
				data := evnt.Data.(map[string]interface{})

				guildId = data["guild"].(string)
				ticket, err = strconv.Atoi(data["ticket"].(string));
				if err != nil {
					conn.Close()
				}

				socket.Guild = guildId
				socket.Ticket = ticket

				// Verify the guild exists
				guildIdParsed, err = strconv.ParseInt(guildId, 10, 64)
				if err != nil {
					fmt.Println(err.Error())
					conn.Close()
					return
				}

				// Get object for selected guild
				var guild objects.Guild
				for _, g := range table.GetGuilds(userIdStr) {
					if g.Id == guildId {
						guild = g
						break
					}
				}

				// Verify the user has permissions to be here
				isAdmin := make(chan bool)
				go utils.IsAdmin(guild, guildIdParsed, userId, isAdmin)
				if !<-isAdmin {
					fmt.Println(err.Error())
					conn.Close()
					return
				}

				// Verify the guild is premium
				premium := make(chan bool)
				go utils.IsPremiumGuild(store, guildId, premium)
				if !<-premium {
					conn.Close()
					return
				}
			} else if evnt.Type == "send" {
				data := evnt.Data.(string)

				if data == "" {
					continue
				}

				// Get ticket UUID from URL and verify it exists
				ticketChan := make(chan table.Ticket)
				go table.GetTicketById(guildIdParsed, ticket, ticketChan)
				ticket := <-ticketChan
				exists := ticket != table.Ticket{}

				contentType := discord.ApplicationJson

				if exists {
					content := data
					if len(content) > 2000 {
						content = content[0:1999]
					}

					// Preferably send via a webhook
					webhookChan := make(chan *string)
					go table.GetWebhookByUuid(ticket.Uuid, webhookChan)
					webhook := <-webhookChan

					success := false
					if webhook != nil {
						success = executeWebhook( ticket.Uuid, *webhook, content, store)
					}

					if !success {
						content = fmt.Sprintf("**%s**: %s", store.Get("name").(string), data)
						if len(content) > 2000 {
							content = content[0:1999]
						}

						endpoint := channel.CreateMessage(int(ticket.Channel))
						err, _ = endpoint.Request(store, &contentType, channel.CreateMessageBody{
							Content: content,
						}, nil)
					}
				}
			}
		}
	}
}

func executeWebhook(uuid, webhook, content string, store sessions.Session) bool {
	body := map[string]interface{}{
		"content":    content,
		"username":   store.Get("name").(string),
		"avatar_url": store.Get("avatar").(string),
	}
	encoded, err := json.Marshal(&body); if err != nil {
		return false
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://canary.discordapp.com/api/webhooks/%s", webhook), bytes.NewBuffer(encoded)); if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = 3 * time.Second

	res, err := client.Do(req); if err != nil {
		return false
	}

	if res.StatusCode == 404 || res.StatusCode == 403 {
		go table.DeleteWebhookByUuid(uuid)
	} else {
		return true
	}

	return false
}
