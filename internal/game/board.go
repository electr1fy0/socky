package game

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const (
	writeWait = time.Second * 5
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

func (b *Board) GenerateFood() {
	x := rand.IntN(b.Rows)
	y := rand.IntN(b.Cols)

	b.mu.Lock()
	b.Grid[x][y] = "f"
	b.mu.Unlock()
}

func (b *Board) Init(rows, cols int) {
	b.Rows, b.Cols = rows, cols
	b.Clients = []*Client{}
	b.Grid = make([][]string, rows)
	for i := range b.Grid {
		b.Grid[i] = make([]string, cols)
		for j := range b.Grid[i] {
			b.Grid[i][j] = "·"
		}
	}
}

func (b *Board) InsertSnake(s *Snake) {
	snakelength := 0
	i := b.Rows / 2
	for j := 2; j < b.Cols; j++ {
		s.Body[snakelength] = Point{X: i + b.SnakeCount, Y: j}
		if snakelength == 0 {
			s.Tail = Point{i + b.SnakeCount, j}
		}
		snakelength++
		if snakelength == MinSnakeLength {
			s.Head = Point{i + b.SnakeCount, j}
			b.SnakeCount++
			return
		}
	}
}

func (b *Board) ClearBoard() {
	if len(b.Clients) == 0 {
		for i := 0; i < b.Rows; i++ {
			for j := 0; j < b.Cols; j++ {
				b.Grid[i][j] = "·"
			}
		}
	}
}

func (b *Board) Update() {
	var toRemove []*Client
	for _, c := range b.Clients {
		c.Snake.Move()
		added := false
		for _, other := range b.Clients {
			if other == c {
				continue
			}
			if slices.Contains(other.Snake.Body, c.Snake.Head) {
				toRemove = append(toRemove, c)
				added = true
			}
		}
		if added {
			continue
		}
		if c.Snake.Head.X < 0 || c.Snake.Head.X >= b.Rows || c.Snake.Head.Y < 0 || c.Snake.Head.Y >= b.Cols {
			toRemove = append(toRemove, c)
			continue
		}

		prevHead := c.Snake.Body[len(c.Snake.Body)-2]
		b.Grid[prevHead.X][prevHead.Y] = "b"
		b.Grid[c.Snake.Tail.X][c.Snake.Tail.Y] = "·"

		if b.Grid[c.Snake.Head.X][c.Snake.Head.Y] == "f" {
			c.Snake.Score++

			c.Snake.Body = append([]Point{c.Snake.Tail}, c.Snake.Body...)
		}
		b.Grid[c.Snake.Head.X][c.Snake.Head.Y] = "h"
	}

	for _, c := range toRemove {
		var over = Message{"over", b.Grid, b.Clients}
		writeCtx, cancel := context.WithTimeout(context.Background(), writeWait)
		wsjson.Write(writeCtx, c.Conn, over)
		cancel()
		b.RemoveClient(c)
		c.Conn.Close(websocket.StatusNormalClosure, "snake died")
	}
}

func (b *Board) ToJSON() ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return json.MarshalIndent(b, "", "  ")
}

func Clear() {
	fmt.Printf("\033[H\033[2J")
}

func (b *Board) AddClient(client *Client) {
	client.Color = snakeColors[(len(b.Clients)+1)%len(snakeColors)]
	b.mu.Lock()
	b.Clients = append(b.Clients, client)
	b.InsertSnake(&client.Snake)

	b.mu.Unlock()
}

func (b *Board) BroadCast() {
	b.mu.RLock()
	clients := make([]*Client, len(b.Clients))

	copy(clients, b.Clients)
	b.mu.RUnlock()
	scores := make(map[string]int)
	colors := make(map[string]string)

	msg := Message{"normal", b.Grid, clients}
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
		writeCtx, cancel := context.WithTimeout(context.Background(), writeWait)
		if err := client.Conn.Write(writeCtx, websocket.MessageText, data); err != nil {
			cancel()
			b.RemoveClient(client)
			client.Conn.Close(websocket.StatusAbnormalClosure, "slow client")
			return
		}
		cancel()
	}
}

func (b *Board) RemoveClient(client *Client) {
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
