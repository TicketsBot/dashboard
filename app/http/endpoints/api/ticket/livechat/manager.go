package livechat

import (
	"encoding/json"
	"github.com/TicketsBot/common/chatrelay"
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
