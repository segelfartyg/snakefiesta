package main

import "fmt"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	spawnClients map[*Client]Chunk

	drunkClients map[*Client]Chunk

	// Inbound messages from the clients.
	broadcast chan []byte
	// Inbound messages from the clients.
	broadcastSpawn chan []byte
	// Inbound messages from the clients.
	broadcastDrunk chan []byte

	chunkChange chan map[string]Chunk

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:      make(chan []byte),
		broadcastSpawn: make(chan []byte),
		broadcastDrunk: make(chan []byte),
		chunkChange:    make(chan map[string]Chunk),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
		spawnClients:   make(map[*Client]Chunk),
		drunkClients:   make(map[*Client]Chunk),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			h.spawnClients[c] = Spawn
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}

		// BROADCASTING TO EVERY CLIENT
		case message := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- message:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}

		// BROADCASTING TO EVERY CLIENT IN THE SPAWN CHUNK
		case spawnMsg := <-h.broadcastSpawn:
			fmt.Println("BROADCASTING SPAWN:", string(spawnMsg))
			for c := range h.spawnClients {
				select {
				case c.send <- spawnMsg:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}

		// BROADCASTING TO EVERY CLIENT IN THE DRUNK CHUNK
		case drunkMsg := <-h.broadcastDrunk:
			fmt.Println("BROADCASTING DRUNK:", string(drunkMsg))
			for c := range h.drunkClients {
				select {
				case c.send <- drunkMsg:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
		// HANDLING A CHUNK CHANGE, RETRIEVING PLAYERID OF THE CHANGE AND ITS NEW CHUNK
		case chunkChanged := <-h.chunkChange:
			var playerId = ""
			var playerChunk = Spawn

			// RETRIEVING THE VALUES
			for id, chunk := range chunkChanged {
				playerId = id
				playerChunk = chunk
			}

			// CHANGING CHUNK ON THE CLIENT LEVEL
			for c := range h.clients {
				if c.playerId == playerId {
					delete(h.spawnClients, c)
					h.drunkClients[c] = Drunk
					c.chunk = playerChunk
				}
			}
		}

	}

}
