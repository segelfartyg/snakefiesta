package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"
)

var boardHeight int = 50
var boardWidth int = 150

var chunkHeight int = 4

const tickRate = 100

type Direction int64
type Chunk int64

const (
	Right Direction = 0
	Down  Direction = 1
	Left  Direction = 2
	Up    Direction = 3
)

const (
	Spawn  Chunk = 0
	Chunky Chunk = 1
	Dunky  Chunk = 2
	Drunk  Chunk = 3
)

type fiestaTile struct {
	X     int
	Y     int
	Chunk Chunk
}

type fiestaBoard struct {
	Rows []fiestaRows
}

type fiestaRows struct {
	Tiles []fiestaTile
}

type playerMovementIntent struct {
	PlayerId  string    `json:"playerId"`
	X         int       `json:"X"`
	Y         int       `json:"Y"`
	Direction Direction `json:"direction"`
}

var fiestaChunks map[Chunk]map[string]fiestaTile = make(map[Chunk]map[string]fiestaTile)

var playerIntents map[string]Direction = make(map[string]Direction)
var playerChunks map[string]Chunk = make(map[string]Chunk)

func main() {

	fiestaChunks[Spawn] = make(map[string]fiestaTile)
	fiestaChunks[Drunk] = make(map[string]fiestaTile)
	flag.Parse()
	hub := newHub()
	go hub.Run()
	go runGameLoop(hub)

	StartServer(hub)
}

func runGameLoop(hub *Hub) {

	// GAME TICKER
	ticker := time.NewTicker(tickRate * time.Millisecond)

	// GAME LOOP
	for range ticker.C {

		movePlayersBasedOnIntent(hub)
		resSpawn, err := json.Marshal(fiestaChunks[Spawn])
		resDrunk, err := json.Marshal(fiestaChunks[Drunk])

		fmt.Println(string(resSpawn))
		fmt.Println(string(resDrunk))
		if err != nil {
			log.Fatal("ERROR GIVING BACK TO BROWSER (RIP CHARITY)")
		}

		hub.broadcastSpawn <- []byte(resSpawn)
		hub.broadcastDrunk <- []byte(resDrunk)
	}
}

func movePlayersBasedOnIntent(hub *Hub) {
	for player, direction := range playerIntents {
		pos := fiestaChunks[playerChunks[player]][player]
		fmt.Println("MOVING PLAYER FROM:")
		fmt.Println(pos.X, " : ", pos.Y)
		fmt.Println("PLAYER CHUNK:")
		fmt.Println(playerChunks[player])
		switch direction {
		case Up:
			if pos.Y-1 < 0 {

				delete(fiestaChunks[playerChunks[player]], player)
				playerChunks[player] = Drunk

				mapToBeSent := map[string]Chunk{
					player: Drunk,
				}

				hub.chunkChange <- mapToBeSent
				pos.Y = boardHeight - 1
				pos.Chunk = Drunk
			} else {
				pos.Y--
			}
			fiestaChunks[playerChunks[player]][player] = pos
		case Right:
			if pos.X+1 > boardWidth-1 {
				pos.X = 0
			} else {
				pos.X++
			}
			fiestaChunks[playerChunks[player]][player] = pos
		case Down:
			if pos.Y+1 > boardHeight-1 {
				pos.Y = 0
			} else {
				pos.Y++
			}
			fiestaChunks[playerChunks[player]][player] = pos
		case Left:
			if pos.X-1 < 0 {
				pos.X = boardWidth - 1
			} else {
				pos.X--
			}
			fiestaChunks[playerChunks[player]][player] = pos
		}
	}
}
