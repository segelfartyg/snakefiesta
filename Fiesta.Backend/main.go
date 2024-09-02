package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

var boardHeight int = 50
var boardWidth int = 50

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

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {

	newFiestaBoard()

	flag.Parse()
	hub := newHub()
	go hub.run()

	go runGameLoop(hub)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
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

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}

	}
}

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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
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

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		intent := playerMovementIntent{}

		fmt.Print("MESSAGE FROM CLIENT:", string(message))

		err = json.Unmarshal([]byte(string(message)), &intent)

		if err != nil {
			log.Fatal("ERROR UNMARSHALLING")
		}

		if _, exists := playerPositions[intent.PlayerId]; exists {
			playerIntents[intent.PlayerId] = intent.Direction
		} else {
			playerIntents[intent.PlayerId] = Right
			playerPositions[intent.PlayerId] = fiestaTile{
				X: 0,
				Y: 0,
			}
		}

		//intent.calculateNextPlayerPosition()

		c.hub.broadcast <- message
	}
}

func (i playerMovementIntent) calculateNextPlayerPosition() {
	fmt.Println("INTENT IN JSON: ", i)

	tile := fiestaTile{X: i.X, Y: i.Y}

	playerPositions[i.PlayerId] = tile
}

func movePlayersBasedOnIntent() {
	for player, direction := range playerIntents {
		pos := playerPositions[player]
		switch direction {
		case Up:
			pos.Y--
			playerPositions[player] = pos
			break
		case Right:
			pos.X++
			playerPositions[player] = pos
			break
		case Down:
			pos.Y++
			playerPositions[player] = pos
			break
		case Left:
			pos.X--
			playerPositions[player] = pos
			break
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
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
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

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

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
