package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"
)

var boardHeight int = 50
var boardWidth int = 50

const tickRate = 1000

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

var playerPositions map[string]fiestaTile = make(map[string]fiestaTile)
var playerIntents map[string]Direction = make(map[string]Direction)

func newFiestaBoard() *fiestaBoard {
	b := fiestaBoard{
		Rows: make([]fiestaRows, boardHeight),
	}
	fmt.Println("INITIAL BOARD")
	for i := range b.Rows {
		fmt.Println("")
		b.Rows[i].Tiles = make([]fiestaTile, boardWidth)

		for j := range b.Rows[i].Tiles {
			b.Rows[i].Tiles[j] = fiestaTile{X: j, Y: i}
			fmt.Print(i, ":", j, ",")
		}
	}
	return &b
}

func main() {

	newFiestaBoard()

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
		res, err := json.Marshal(playerPositions)

		if err != nil {
			log.Fatal("ERROR GIVING BACK TO BROWSER (RIP CHARITY)")
		}

		hub.broadcast <- []byte(res)
	}
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (i playerMovementIntent) calculateNextPlayerPosition() {
	fmt.Println("INTENT IN JSON: ", i)
	tile := fiestaTile{X: i.X, Y: i.Y}
	playerPositions[i.PlayerId] = tile
}

func movePlayersBasedOnIntent() {
	for player, direction := range playerIntents {
		pos := playerPositions[player]
		fmt.Println("MOVING PLAYER FROM:")
		fmt.Println(pos.X, " : ", pos.Y)
		switch direction {
		case Up:
			if pos.Y-1 < 0 {
				pos.Y = boardHeight - 1
			} else {
				pos.Y--
			}
			playerPositions[player] = pos
		case Right:
			if pos.X+1 > boardWidth-1 {
				pos.X = 0
			} else {
				pos.X++
			}
			playerPositions[player] = pos
		case Down:
			if pos.Y+1 > boardHeight-1 {
				pos.Y = 0
			} else {
				pos.Y++
			}
			playerPositions[player] = pos
		case Left:
			if pos.X-1 < 0 {
				pos.X = boardWidth - 1
			} else {
				pos.X--
			}
			playerPositions[player] = pos
		}
	}
}
