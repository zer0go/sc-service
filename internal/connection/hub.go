package connection

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound message from the clients.
	message chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		message:    make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.message:
			h.forwardMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.clients[client.id] = client

	log.Debug().
		Str("clientId", client.id).
		Msg("client registered")
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client.id]; !ok {
		log.Warn().
			Str("clientId", client.id).
			Msg("unregister client failed, client not connected")
		return
	}

	delete(h.clients, client.id)
	close(client.send)

	log.Debug().
		Str("clientId", client.id).
		Msg("client unregistered")
}

func (h *Hub) forwardMessage(message *Message) {
	if _, ok := h.clients[message.RecipientId]; !ok {
		log.Warn().
			Str("recipientId", message.RecipientId).
			Msg("send message failed, recipient not connected")

		return
	}

	m, err := json.Marshal(message)
	if err != nil {
		log.Err(err).
			Msg("encode json failed")

		return
	}

	h.clients[message.RecipientId].send <- m

	log.Debug().
		Str("recipientId", message.RecipientId).
		Msg("message sent")
}
