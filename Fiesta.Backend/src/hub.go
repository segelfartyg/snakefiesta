package main

import (
	"encoding/json"
	"fmt"
	"log"
)

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

	// Inbound messages from the clients.
	broadcastAll chan map[Chunk]map[string]fiestaTile

	chunkChange chan map[string]Chunk

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	klienter map[Chunk]map[*Client]Chunk
}

func newHub(chunks []Chunk) *Hub {

	klienter := make(map[Chunk]map[*Client]Chunk)

	for _, v := range chunks {
		klienter[v] = make(map[*Client]Chunk)
	}

	fmt.Println(klienter)
	return &Hub{
		klienter:       klienter,
		broadcast:      make(chan []byte),
		broadcastSpawn: make(chan []byte),
		broadcastDrunk: make(chan []byte),
		broadcastAll:   make(chan map[Chunk]map[string]fiestaTile),
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
			h.klienter[Spawn][c] = Spawn
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
		case entireBoard := <-h.broadcastAll:
			fmt.Println("BROADCASTING DRUNK:", entireBoard)

			for key, clientOfBoard := range entireBoard {
				fmt.Println("KEY: ", key)
				fmt.Println("CLIENT: ", clientOfBoard)
			}

			for key, _ := range h.klienter {

				for c := range h.klienter[key] {
					fmt.Print("HEJ DU: ", c.playerId, c.chunk)

					sendJson, err := json.Marshal(entireBoard[key])
					if err != nil {
						log.Fatal("ERROR GIVING BACK TO BROWSER (RIP CHARITY)")
					}

					c.send <- []byte(sendJson)

				}
			}

			// }

			// for c := range h.klienter {
			// 	select {
			// 	case c.send <- klientMsg:
			// 	default:
			// 		close(c.send)
			// 		delete(h.clients, c)
			// 	}
			// }
		// HANDLING A CHUNK CHANGE, RETRIEVING PLAYERID OF THE CHANGE AND ITS NEW CHUNK
		case chunkChanged := <-h.chunkChange:
			var playerId = ""
			var playerChunk = Spawn

			// RETRIEVING THE VALUES
			for id, chunk := range chunkChanged {
				playerId = id
				playerChunk = chunk
			}
			fmt.Println(playerId, playerChunk)
			// CHANGING CHUNK ON THE CLIENT LEVEL
			for c := range h.clients {
				if c.playerId == playerId {
					delete(h.klienter[Spawn], c)
					h.klienter[playerChunk][c] = playerChunk
					c.chunk = playerChunk
				}
			}
		}

	}

}
