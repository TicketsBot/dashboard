package root

import (
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == config.Conf.Server.BaseUrl
	},
}

var SocketsLock sync.RWMutex
var Sockets []*Socket

type (
	Socket struct {
		Ws       *websocket.Conn
		GuildId  uint64
		TicketId int
	}

	WsEvent struct {
		Type string
		Data json.RawMessage
	}

	AuthEvent struct {
		GuildId  uint64 `json:"guild_id,string"`
		TicketId int    `json:"ticket_id"`
		Token    string `json:"token"`
	}
)

var timeout = time.Second * 60

func WebChatWs(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	socket := &Socket{
		Ws: conn,
	}

	SocketsLock.Lock()
	Sockets = append(Sockets, socket)
	SocketsLock.Unlock()

	conn.SetCloseHandler(func(code int, text string) error {
		i := -1
		SocketsLock.Lock()
		defer SocketsLock.Unlock()

		for index, element := range Sockets {
			if element == socket {
				i = index
				break
			}
		}

		if i != -1 {
			Sockets = Sockets[:i+copy(Sockets[i:], Sockets[i+1:])]
		}

		return nil
	})

	lastResponse := time.Now()
	conn.SetPongHandler(func(a string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		// We can let this func call the CloseHandler
		for {
			err := conn.WriteMessage(websocket.PingMessage, []byte("keepalive"))
			if err != nil {
				conn.Close()
				conn.CloseHandler()(1000, "")
				return
			}

			time.Sleep(timeout / 2)
			if time.Since(lastResponse) > timeout {
				conn.Close()
				conn.CloseHandler()(1000, "")
				return
			}
		}
	}()

	for {
		var event WsEvent
		err := conn.ReadJSON(&event)
		if err != nil {
			break
		}

		if socket.GuildId == 0 && event.Type != "auth" {
			conn.Close()
			break
		} else if event.Type == "auth" {
			var authData AuthEvent
			if err := json.Unmarshal(event.Data, &authData); err != nil {
				conn.Close()
				return
			}

			token, err := jwt.Parse(authData.Token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(config.Conf.Server.Secret), nil
			})

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				conn.Close()
				return
			}

			userIdStr, ok := claims["userid"].(string)
			if !ok {
				conn.Close()
				return
			}

			userId, err := strconv.ParseUint(userIdStr, 10, 64)
			if err != nil {
				conn.Close()
				return
			}

			// Verify the user has permissions to be here
			permLevel, err := utils.GetPermissionLevel(authData.GuildId, userId)
			if err != nil {
				conn.Close()
				return
			}

			if permLevel < permission.Admin {
				conn.Close()
				return
			}

			botContext, err := botcontext.ContextForGuild(authData.GuildId)
			if err != nil {
				ctx.JSON(500, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}

			// Verify the guild is premium
			premiumTier, err := rpc.PremiumClient.GetTierByGuildId(authData.GuildId, true, botContext.Token, botContext.RateLimiter)
			if err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}

			if premiumTier == premium.None {
				conn.Close()
				return
			}

			SocketsLock.Lock()
			socket.GuildId = authData.GuildId
			socket.TicketId = authData.TicketId
			SocketsLock.Unlock()
		}
	}
}
