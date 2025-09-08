package board

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// TODO: make members that shouldn't be exported lower case again

const (
	Min        = 5
	TickRate   = 300 * time.Millisecond
	FoodPeriod = 4 * time.Second
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Point struct{ X, Y int }

type Snake struct {
	Body      []Point
	Head      Point
	Tail      Point
	Size      int
	Direction Direction
	Score     int
}

func (s *Snake) Init() {
	s.Body = make([]Point, Min+1)
	s.Direction = Right
	s.Score = 0
}

func (s *Snake) shift(newHead Point) {
	n := len(s.Body)
	s.Tail = Point{s.Body[0].X, s.Body[0].Y}
	for i := 0; i < n-1; i++ {
		s.Body[i].X = s.Body[i+1].X
		s.Body[i].Y = s.Body[i+1].Y
	}
	s.Body[n-1] = newHead
	s.Head = newHead

}

func (s *Snake) Move() {
	switch s.Direction {
	case Up:
		newHead := Point{s.Head.X - 1, s.Head.Y}
		s.shift(newHead)
	case Down:
		newHead := Point{s.Head.X + 1, s.Head.Y}
		s.shift(newHead)
	case Left:
		newHead := Point{s.Head.X, s.Head.Y - 1}
		s.shift(newHead)
	case Right:
		newHead := Point{s.Head.X, s.Head.Y + 1}
		s.shift(newHead)
	}
}

type Board struct {
	Rows, Cols int
	Food       Point
	Snakes     []*Snake
	Grid       [][]rune
	SnakeCount int
}

func (b *Board) GenerateFood() {
	time.Sleep(4 * time.Second)
	x := rand.IntN(b.Rows)
	y := rand.IntN(b.Cols)

	b.Food = Point{X: x, Y: y}
	b.Grid[x][y] = '⊗'
}

func (b *Board) Init(rows, cols int) {
	b.Rows, b.Cols = rows, cols
	b.Snakes = []*Snake{}
	b.Grid = make([][]rune, rows)
	for i := range b.Grid {
		b.Grid[i] = make([]rune, cols)
		for j := range b.Grid[i] {
			b.Grid[i][j] = '.'
		}
	}
	b.SnakeCount = 0
	go b.GenerateFood()
}

func (b *Board) InitSnake(s *Snake) {
	snakelength := 0
	i := b.Rows / 2
	for j := 2; j < b.Cols; j++ {
		s.Body[snakelength] = Point{X: i + b.SnakeCount, Y: j}
		if snakelength == 0 {
			s.Tail = Point{i + b.SnakeCount, j}
		}
		snakelength++
		if snakelength == Min {
			s.Head = Point{i + b.SnakeCount, j}
			b.SnakeCount++
			return
		}
	}
}

func (b *Board) Update() {
	for _, s := range b.Snakes {
		s.Move()

		if s.Head == b.Food {
			s.Score++
			b.Grid[s.Head.X][s.Head.Y] = '.'
			s.Body = append([]Point{s.Tail}, s.Body...)
			// REMINDER: add tail logic
			go b.GenerateFood()
		}
	}
}

func (b *Board) Print() {
	var output strings.Builder

	for i := 0; i < b.Rows; i++ {
		for j := 0; j < b.Cols; j++ {
			printed := false
			if b.SnakeCount != 0 {
				for _, snake := range b.Snakes {
					for _, point := range snake.Body {
						if point.X == i && point.Y == j {
							output.WriteString("◉ ")
							printed = true
						}
					}
				}
			}
			if !printed {
				output.WriteString(string(b.Grid[i][j]) + " ")
			}
		}
		output.WriteString("\n\r")
	}

	fmt.Print(output.String())
}

func Clear() {
	fmt.Printf("\033[H\033[2J")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	Conn     *websocket.Conn
	ID       string
	Keypress string
	Snake    Snake
}

var (
	clients   = make(map[*Client]bool)
	clientsMu sync.Mutex
)

func (b *Board) Run(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	client := &Client{
		Conn: conn, ID: r.RemoteAddr,
	}
	client.Snake.Init()

	clientsMu.Lock()
	b.Snakes = append(b.Snakes, &client.Snake)
	b.InitSnake(&client.Snake)
	clientsMu.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Msg reading err:", err)
			break
		}

		client.Keypress = string(msg)

		switch client.Keypress {
		case "k":
			if client.Snake.Direction != Down {
				client.Snake.Direction = Up
			}
		case "j":
			if client.Snake.Direction != Up {
				client.Snake.Direction = Down
			}
		case "l":
			if client.Snake.Direction != Left {
				client.Snake.Direction = Right
			}
		case "h":
			if client.Snake.Direction != Right {
				client.Snake.Direction = Left
			}
		default:
			continue
		}

		fmt.Println("received:", client.Keypress)
	}
}
