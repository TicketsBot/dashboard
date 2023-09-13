package livechat

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	Manager       *SocketManager
	Ws            *websocket.Conn
	RequestCtx    *gin.Context
	Authenticated bool
	GuildId       uint64
	TicketId      int
	tx            chan any
	flush         chan chan struct{}
}

const (
	messageSizeLimit   = 1024 * 32
	keepaliveFrequency = 45 * time.Second
	keepaliveTimeout   = 60 * time.Second
	writeTimeout       = 10 * time.Second
)

func NewClient(manager *SocketManager, ws *websocket.Conn, c *gin.Context, guildId uint64, ticketId int) *Client {
	return &Client{
		Manager:       manager,
		Ws:            ws,
		RequestCtx:    c,
		Authenticated: false,
		GuildId:       guildId,
		TicketId:      ticketId,
		tx:            make(chan any),
		flush:         make(chan chan struct{}),
	}
}

func (c *Client) Close() {
	close(c.tx)
}

func (c *Client) StartReadLoop() error {
	defer func() {
		c.Manager.unregister <- c
		_ = c.Ws.Close()
		c.Close()
	}()

	// Set up connection properties
	c.Ws.SetReadLimit(messageSizeLimit)
	if err := c.Ws.SetReadDeadline(time.Now().Add(keepaliveTimeout)); err != nil {
		return err
	}

	c.Ws.SetPongHandler(func(appData string) error {
		return c.Ws.SetReadDeadline(time.Now().Add(keepaliveTimeout))
	})

	for {
		var event Event
		if err := c.Ws.ReadJSON(&event); err != nil {
			return err
		}

		if !c.Authenticated && event.Type != EventTypeAuth {
			if err := c.Ws.WriteJSON(NewErrorMessage("Unauthorized")); err != nil {
				return err
			}

			return nil
		}

		if err := c.HandleEvent(event); err != nil {
			c.RequestCtx.Error(err)
			c.Write(NewErrorMessage(err.Error()))
			c.Flush()
			_ = c.Ws.Close()
			return err
		}
	}
}

func (c *Client) Write(msg any) {
	c.tx <- msg
}

func (c *Client) StartWriteLoop() error {
	ticker := time.NewTicker(keepaliveFrequency)
	defer func() {
		ticker.Stop()
		_ = c.Ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.tx:
			if err := c.Ws.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
				return err
			}

			if !ok { // Channel was closed
				_ = c.Ws.WriteMessage(websocket.CloseMessage, []byte{})
				return nil
			} else {
				if err := c.Ws.WriteJSON(message); err != nil {
					return err
				}
			}
		case <-ticker.C:
			if err := c.Ws.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
				return err
			}

			if err := c.Ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		case ch := <-c.flush: // TODO: Channel order is random, there is a race condition here
			ch <- struct{}{}
		}
	}
}

func (c *Client) Flush() {
	ch := make(chan struct{})
	c.flush <- ch

	timer := time.After(time.Second)
	select {
	case <-ch:
		return
	case <-timer:
		return
	}
}
