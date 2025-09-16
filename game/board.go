package game

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Board struct {
	Rows        int          `json:"rows"`
	Cols        int          `json:"cols"`
	Grid        [][]string   `json:"grid"`
	GridString  string       `json:"gridString"`
	SnakeCount  int          `json:"snakeCount"`
	Clients     []*Client    `json:"clients"`
	ClientCount int          `json:"-"`
	mu          sync.RWMutex `json:"-"`
}

type Client struct {
	ID       string          `json:"id"`
	Keypress string          `json:"keypress"`
	Snake    Snake           `json:"snake"`
	Name     string          `json:"name"`
	Color    string          `json:"color"`
	Conn     *websocket.Conn `json:"-"`
}
type Message struct {
	Type    string     `json:"type"`
	Grid    [][]string `json:"grid"`
	Clients []*Client  `json:"clients"`
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
		c.Conn.WriteJSON(over)
		b.removeClient(c)
		c.Conn.Close()
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
