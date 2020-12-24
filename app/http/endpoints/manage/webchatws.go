package manage

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/permission"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strconv"
	"sync"
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
	userId := utils.GetUserId(store)

	var guildId string
	var guildIdParsed uint64
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
			ticket, err = strconv.Atoi(data["ticket"].(string))
			if err != nil {
				conn.Close()
				break
			}

			socket.Guild = guildId
			socket.Ticket = ticket

			// Verify the guild exists
			guildIdParsed, err = strconv.ParseUint(guildId, 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				conn.Close()
				return
			}

			// Verify the user has permissions to be here
			permLevel, err := utils.GetPermissionLevel(guildIdParsed, userId)
			if err != nil {
				fmt.Println(err.Error())
				conn.Close()
				return
			}

			if permLevel < permission.Admin {
				fmt.Println(err.Error())
				conn.Close()
				return
			}

			botContext, err := botcontext.ContextForGuild(guildIdParsed)
			if err != nil {
				ctx.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"error": err.Error(),
				})
				return
			}

			// Verify the guild is premium
			premiumTier := rpc.PremiumClient.GetTierByGuildId(guildIdParsed, true, botContext.Token, botContext.RateLimiter)
			if premiumTier == premium.None {
				conn.Close()
				return
			}
		}
	}
}
