// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package connection

import (
	"github.com/rs/zerolog/log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 8096
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub  *Hub
	id   string
	conn *websocket.Conn
	send chan []byte
}

func NewClient(hub *Hub, id string, conn *websocket.Conn) *Client {
	client := &Client{
		hub:  hub,
		id:   id,
		conn: conn,
		send: make(chan []byte, 256),
	}
	go client.writePump()
	go client.readPump()

	return client
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		var message Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			log.Err(err).Msg("invalid message json received")

			break
		}

		if message.RecipientId == "" {
			log.Warn().Msg("empty recipient id received")
			continue
		}

		log.Debug().
			Str("senderId", c.id).
			Str("recipientId", message.RecipientId).
			Msg("receive message")

		c.hub.message <- &message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Err(err).Msg("next writer error")
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Err(err).Msg("write message failed")
				return
			}

			if err = w.Close(); err != nil {
				log.Err(err).Msg("close writer failed")
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
