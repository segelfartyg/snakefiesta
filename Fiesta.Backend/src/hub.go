package main

import (
	"encoding/json"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// CHANNEL FOR BROADCASTING TO ALL CLIENTS
	broadcastAll chan map[Chunk]map[string]fiestaTile

	// CHANNEL FOR LISTENING TO CHANGE OF CHUNK
	chunkChange chan chunkChangeEvent

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	chunkClients map[Chunk]map[*Client]bool
}

func newHub(chunks []Chunk) *Hub {

	chunkClients := make(map[Chunk]map[*Client]bool)

	for _, v := range chunks {
		chunkClients[v] = make(map[*Client]bool)
	}

	return &Hub{
		chunkClients: chunkClients,
		broadcast:    make(chan []byte),
		broadcastAll: make(chan map[Chunk]map[string]fiestaTile),
		chunkChange:  make(chan chunkChangeEvent),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			h.chunkClients[Spawn][c] = true
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

		// GETTING EVENT WITH THE ENTIRE BOARD
		case entireBoard := <-h.broadcastAll:

			// LOOPING THROUGH EVERY CLIENT FOR EACH CHUNK
			for chunk, _ := range h.chunkClients {

				// LOOPING THROUGH EACH CLIENT IN THE CHUNK
				for clientInChunk := range h.chunkClients[chunk] {
					sendJson, err := json.Marshal(entireBoard[chunk])
					if err != nil {
						log.Fatal("ERROR GIVING BACK TO BROWSER (RIP CHARITY)")
					}
					// SENDING CHUNK DATA TO CLIENTS OF THE CHUNK
					clientInChunk.send <- []byte(sendJson)
				}
			}

		// HANDLING A CHUNK CHANGE, RETRIEVING PLAYERID OF THE CHANGE AND ITS NEW CHUNK
		case chunkChange := <-h.chunkChange:

			// RETRIEVING THE VALUES
			playerId := chunkChange.PlayerId
			previousChunk := chunkChange.PreviousChunk
			nextChunk := chunkChange.NextChunk

			// CHANGING CHUNK ON THE CONNECTION LEVEL
			for c := range h.clients {
				if c.playerId == playerId {
					delete(h.chunkClients[previousChunk], c)
					h.chunkClients[nextChunk][c] = true
					c.chunk = nextChunk
				}
			}
		}

	}

}
