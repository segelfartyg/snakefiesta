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

const tickRate = 100

type Direction int64

const (
	Right Direction = 0
	Down  Direction = 1
	Left  Direction = 2
	Up    Direction = 3
)

type fiestaTile struct {
	X int
	Y int
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

var fiestaChunkList = [2]string{"spawnChunk", "chunkyCunk"}
var fiestaChunks map[string]map[string]fiestaTile = make(map[string]map[string]fiestaTile)
var playerIntents map[string]Direction = make(map[string]Direction)

func main() {

	fiestaChunks["spawnChunk"] = make(map[string]fiestaTile)
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

		movePlayersBasedOnIntent()
		res, err := json.Marshal(fiestaChunks["spawnChunk"])

		if err != nil {
			log.Fatal("ERROR GIVING BACK TO BROWSER (RIP CHARITY)")
		}

		hub.broadcast <- []byte(res)
	}
}

func movePlayersBasedOnIntent() {
	for player, direction := range playerIntents {
		pos := fiestaChunks["spawnChunk"][player]
		fmt.Println("MOVING PLAYER FROM:")
		fmt.Println(pos.X, " : ", pos.Y)
		switch direction {
		case Up:
			if pos.Y-1 < 0 {
				pos.Y = boardHeight - 1
			} else {
				pos.Y--
			}
			fiestaChunks["spawnChunk"][player] = pos
		case Right:
			if pos.X+1 > boardWidth-1 {
				pos.X = 0
			} else {
				pos.X++
			}
			fiestaChunks["spawnChunk"][player] = pos
		case Down:
			if pos.Y+1 > boardHeight-1 {
				pos.Y = 0
			} else {
				pos.Y++
			}
			fiestaChunks["spawnChunk"][player] = pos
		case Left:
			if pos.X-1 < 0 {
				pos.X = boardWidth - 1
			} else {
				pos.X--
			}
			fiestaChunks["spawnChunk"][player] = pos
		}
	}
}
