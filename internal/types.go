package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
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

const (
	boardHeight = 30
	boardWidth  = 35
	FoodPeriod  = 8 * time.Second
	TickRate    = 150 * time.Millisecond
)

type Room struct {
	ID string

	Board  *Board
	ticker *time.Ticker
	done   chan struct{}
}

type RoomManager struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

func NewRoom() *Room {
	b := &Board{}
	b.Init(boardHeight, boardWidth)
	id := uuid.NewString()
	fmt.Println("Room ID: ", id)

	return &Room{
		ID:    id,
		Board: b,
		done:  make(chan struct{}),
	}
}

func (R *Room) InitRoom() {
	tick := time.Tick(TickRate)
	foodTick := time.Tick(FoodPeriod)
	for {
		select {
		case <-tick:
			R.Board.Update()
			R.Board.ClearBoard()
			R.Board.BroadCast()
		case <-foodTick:
			R.Board.GenerateFood()
		case <-R.done:
			return
		}
	}
}

func (R *Room) Close() {
	select {
	case <-R.done:
	default:
		close(R.done)
	}
}
