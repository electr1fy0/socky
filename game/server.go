package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var snakeColors = []string{
	"red",
	"yellow",
	"green",
	"blue",
	"magenta",
	"cyan",
	"white",
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (b *Board) addClient(client *Client) {
	b.ClientCount++
	client.Color = snakeColors[b.ClientCount%len(snakeColors)]
	b.mu.Lock()
	b.Clients = append(b.Clients, client)
	b.InsertSnake(&client.Snake)

	b.mu.Unlock()
}

func (b *Board) BroadCast() {
	// b.Print()
	b.mu.RLock()
	// boardState := b.Print()
	clients := make([]*Client, len(b.Clients))
	// gridString := b.GridString
	copy(clients, b.Clients)
	b.mu.RUnlock()
	scores := make(map[string]int)
	colors := make(map[string]string)

	msg := Message{b.Grid, clients}
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error marshalling:", err)
		return
	}
	for _, client := range clients {
		scores[client.Name] = client.Snake.Score
		colors[client.Name] = client.Color
	}
	for _, client := range clients {

		if err := client.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			b.removeClient(client)
			client.Conn.Close()
		}
	}
}

func (b *Board) removeClient(client *Client) {
	b.mu.Lock()

	for i, c := range b.Clients {
		if c.ID == client.ID {
			for _, point := range c.Snake.Body {
				if point.X >= 0 && point.X < b.Rows && point.Y >= 0 && point.Y < b.Cols {
					b.Grid[point.X][point.Y] = "·"
				}
			}
			b.Grid[c.Snake.Tail.X][c.Snake.Tail.Y] = "·"
			b.Clients = append(b.Clients[:i], b.Clients[i+1:]...)
			b.SnakeCount--
			break
		}
	}
	b.mu.Unlock()
}

func getKeypresses(client *Client) {
	for {
		_, msg, err := client.Conn.ReadMessage()

		if err != nil {
			fmt.Println("Err reading:", err)
			break
		}
		msgString := string(msg)
		if name, ok := strings.CutPrefix(msgString, "NAME:"); ok {
			client.Name = name
			continue
		}

		client.Keypress = string(msg)

		switch client.Keypress {
		case "k", "w":
			if client.Snake.Direction != Down {
				client.Snake.Direction = Up
			}
		case "j", "s":
			if client.Snake.Direction != Up {
				client.Snake.Direction = Down
			}
		case "l", "d":
			if client.Snake.Direction != Left {
				client.Snake.Direction = Right
			}
		case "h", "a":
			if client.Snake.Direction != Right {
				client.Snake.Direction = Left
			}
		default:
			continue
		}
		fmt.Print("received:", client.Keypress, "\r")
	}
}

func (b *Board) Run(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	client := &Client{
		Conn: conn, ID: r.RemoteAddr,
	}
	defer func() {
		b.removeClient(client)
		fmt.Println("client left")
		conn.Close()
	}()
	client.Snake.Init()
	b.addClient(client)

	getKeypresses(client)
}
