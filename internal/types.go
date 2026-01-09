package internal

import (
	"sync"

	"github.com/coder/websocket"
)

type Direction int

const MinSnakeLength = 5

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Snake struct {
	Body      []Point
	Head      Point
	Tail      Point
	Direction Direction
	Score     int
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Board struct {
	Rows       int          `json:"rows"`
	Cols       int          `json:"cols"`
	Grid       [][]string   `json:"grid"`
	GridString string       `json:"gridString"`
	SnakeCount int          `json:"snakeCount"`
	Clients    []*Client    `json:"clients"`
	mu         sync.RWMutex `json:"-"`
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