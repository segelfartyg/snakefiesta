package main

import (
	"encoding/json"
	"log"
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
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	chunk Chunk

	playerId string
}

// READS MESSAGES COMING FROM THE CLIENT
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		intent := playerMovementIntent{}

		err = json.Unmarshal([]byte(string(message)), &intent)

		if err != nil {
			log.Fatal("ERROR UNMARSHALLING")
		}

		// TODO: SET PLAYER ID DURING WEBSOCKET CONNECTION REGISTRATION
		c.playerId = intent.PlayerId

		// SETTING INLY INTENT IF THERE IS A PLAYER PRESENT ALREADY
		if _, exists := playerIntents[intent.PlayerId]; exists {
			playerIntents[intent.PlayerId] = intent.Direction
		} else {
			// REGISTERING THE CLIENT IN THE SPAWN CHUNK.
			playerIntents[intent.PlayerId] = Right
			playerChunkPositions[intent.PlayerId] = getSpawnChunk()
			fiestaChunks[Spawn][intent.PlayerId] = fiestaTile{
				X:     0,
				Y:     0,
				Chunk: Spawn,
			}
		}
	}
}

// SENDS MESSAGES TO THE CLIENT
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)

			if err != nil {
				log.Fatal("ERROR MARSHALLING TO CLIENT (RIP?)")
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
