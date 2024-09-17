package main

import (
	"flag"
	"fmt"
	"time"
)

var boardHeight int = 50
var boardWidth int = 150

const tickRate = 50

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
	Drunk  Chunk = 1
	Dunky  Chunk = 2
	Chunky Chunk = 3
	Lunky  Chunk = 4
	Sunky  Chunk = 5
)

type chunkChangeEvent struct {
	PreviousChunk Chunk
	NextChunk     Chunk
	PlayerId      string
}

type fiestaTile struct {
	X     int
	Y     int
	Chunk Chunk
}

type playerMovementIntent struct {
	PlayerId  string    `json:"playerId"`
	X         int       `json:"X"`
	Y         int       `json:"Y"`
	Direction Direction `json:"direction"`
}

var chunks = []Chunk{Spawn, Drunk, Dunky, Chunky, Lunky, Sunky}

// AN ENTIRE BOARD, EVERY CHUNK MAPS TO ANOTHER MAP THAT MAPS PLAYERID TO A POSITION IN THAT CHUNK
var fiestaChunks map[Chunk]map[string]fiestaTile = make(map[Chunk]map[string]fiestaTile)

// MAP FOR STORING WHERE EVERY PLAYER WISHES TO MOVE
var playerIntents map[string]Direction = make(map[string]Direction)

// MAP FOR STORING WHERE EVERY PLAYER IS
var playerChunkPositions map[string]Chunk = make(map[string]Chunk)

func main() {

	for _, c := range chunks {
		fiestaChunks[c] = make(map[string]fiestaTile)
	}

	flag.Parse()
	hub := newHub(chunks)

	go hub.Run()
	go runGameLoop(hub)

	StartServer(hub)
}

func runGameLoop(hub *Hub) {

	// GAME TICKER
	ticker := time.NewTicker(tickRate * time.Millisecond)

	// GAME LOOP
	for range ticker.C {

		// MANIPULATING THE fiestaChunk BOARD AND ITS STATE
		movePlayersBasedOnIntent(hub)

		// BROADCASTING BOARD TO EACH CHUNK AND THEIR CLIENTS
		hub.broadcastAll <- fiestaChunks
	}
}

func movePlayersBasedOnIntent(hub *Hub) {
	for player, direction := range playerIntents {

		// RETRIEVING THE POSITION OF THE PLAYER, IN THE CHUNK OF THE PLAYER
		pos := fiestaChunks[playerChunkPositions[player]][player]
		switch direction {
		case Up:
			if pos.Y-1 < 0 {
				pos.Y = boardHeight - 1
			} else {
				pos.Y--
			}
			fiestaChunks[playerChunkPositions[player]][player] = pos
		case Right:
			if pos.X+1 > boardWidth-1 {
				pos.X = 0
				handleChunkExit(direction, player, &pos, hub)
			} else {
				pos.X++
			}
			fiestaChunks[playerChunkPositions[player]][player] = pos
		case Down:
			if pos.Y+1 > boardHeight-1 {
				pos.Y = 0
			} else {
				pos.Y++
			}
			fiestaChunks[playerChunkPositions[player]][player] = pos
		case Left:
			if pos.X-1 < 0 {
				pos.X = boardWidth - 1
				handleChunkExit(direction, player, &pos, hub)
			} else {
				pos.X--
			}
			fiestaChunks[playerChunkPositions[player]][player] = pos
		}
	}
}

func getSpawnChunk() Chunk {
	return chunks[0]
}

func handleChunkExit(direction Direction, playerId string, position *fiestaTile, hub *Hub) {
	delete(fiestaChunks[playerChunkPositions[playerId]], playerId)

	nextChunk := getNextChunk(direction, playerId)
	changeEvent := chunkChangeEvent{
		PreviousChunk: playerChunkPositions[playerId],
		NextChunk:     nextChunk,
		PlayerId:      playerId,
	}

	playerChunkPositions[playerId] = nextChunk
	position.Chunk = nextChunk
	fmt.Println("NEXT CHUNK", nextChunk)
	fmt.Println("PREV CHUNK", changeEvent.PreviousChunk)
	hub.chunkChange <- changeEvent
}

func getNextChunk(direction Direction, playerId string) Chunk {

	switch direction {
	case Up:
		return playerChunkPositions[playerId]
	case Right:
		return chunks[playerChunkPositions[playerId]+1]
	case Down:
		return playerChunkPositions[playerId]
	case Left:
		return chunks[playerChunkPositions[playerId]-1]
	}

	return Spawn
}
