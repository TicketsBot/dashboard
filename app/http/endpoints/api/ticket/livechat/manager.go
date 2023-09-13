package livechat

import (
	"encoding/json"
	"github.com/TicketsBot/common/chatrelay"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
)

var (
	activeWebsockets = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "tickets",
		Subsystem: "api",
		Name:      "active_livechat_websockets",
		Help:      "The number of open live-chat websockets",
	})

	websocketMessages = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "tickets",
		Subsystem: "api",
		Name:      "livechat_websocket_messages",
		Help:      "The number of messages relayed over live-chat websockets",
	}, []string{"guild_id", "message_id"})
)

type (
	SocketManager struct {
		clients    map[uint64][]*Client // Remember: A client might not be authenticated!
		messages   chan chatrelay.MessageData
		register   chan *Client
		unregister chan *Client
	}
)

func NewSocketManager() *SocketManager {
	return &SocketManager{
		clients:    map[uint64][]*Client{},
		messages:   make(chan chatrelay.MessageData),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (sm *SocketManager) Run() {
	for {
		select {
		case client := <-sm.register:
			guildClients := sm.clients[client.GuildId]
			guildClients = append(guildClients, client)
			sm.clients[client.GuildId] = guildClients

			activeWebsockets.Inc()
		case client := <-sm.unregister:
			guildClients := sm.clients[client.GuildId]
			if len(guildClients) == 0 {
				continue // TODO: Warn
			}

			i := -1
			for index, el := range guildClients {
				if el == client {
					i = index
					break
				}
			}

			if i != -1 {
				guildClients = guildClients[:i+copy(guildClients[i:], guildClients[i+1:])]
			}

			sm.clients[client.GuildId] = guildClients

			activeWebsockets.Dec()
		case msg := <-sm.messages:
			guildClients, ok := sm.clients[msg.Ticket.GuildId]
			if !ok || len(guildClients) == 0 { // No clients connected to this API server for this guild
				continue
			}

			encoded, err := json.Marshal(msg.Message)
			if err != nil {
				continue // TODO: Warn
			}

			for _, client := range guildClients {
				if !client.Authenticated {
					continue
				}

				// Should already be filtered by guild ID, but here we are filtering by ticket ID for the first time
				if client.GuildId != msg.Ticket.GuildId || client.TicketId != msg.Ticket.Id {
					continue
				}

				websocketMessages.WithLabelValues(
					strconv.FormatUint(client.GuildId, 10),
					strconv.FormatUint(msg.Message.Id, 10),
				).Inc()

				client.Write(Event{
					Type: EventTypeMessage,
					Data: encoded,
				})
			}
		}
	}
}

func (sm *SocketManager) BroadcastMessage(message chatrelay.MessageData) {
	sm.messages <- message
}
