package main

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

const (
	MIN        = 5
	tickRate   = 300 * time.Millisecond
	foodPeriod = 4 * time.Second
)

type Direction int

const (
	up Direction = iota
	down
	left
	right
)

type Point struct{ x, y int }

// snake stuff
type Snake struct {
	body      []Point
	head      Point
	tail      Point
	size      int
	direction Direction
	score     int
}

func (s *Snake) init() {
	s.body = make([]Point, MIN+1)
	s.direction = right

	s.score = 0
}

func (s *Snake) shift(newHead Point) {
	n := len(s.body)
	for i := 0; i < n-1; i++ {
		s.body[i].x = s.body[i+1].x
		s.body[i].y = s.body[i+1].y
	}
	s.body[n-1] = newHead
	s.head = newHead
}

func (s *Snake) move() {
	switch s.direction {
	case up:
		newHead := Point{s.head.x - 1, s.head.y}
		s.shift(newHead)
	case down:
		newHead := Point{s.head.x + 1, s.head.y}
		s.shift(newHead)
	case left:
		newHead := Point{s.head.x, s.head.y - 1}
		s.shift(newHead)
	case right:
		newHead := Point{s.head.x, s.head.y + 1}
		s.shift(newHead)
	}
}

// board stuff
type Board struct {
	rows, cols int
	food       Point
	snakes     []*Snake
	grid       [][]rune
	snakeCount int
}

func (b *Board) generateFood() {
	time.Sleep(4 * time.Second)
	x := rand.IntN(b.rows)
	y := rand.IntN(b.cols)

	b.food = Point{x, y}

	b.grid[x][y] = '⊗'
}

func (b *Board) init(rows, cols int) {
	b.rows, b.cols = rows, cols
	b.snakes = []*Snake{}
	b.grid = make([][]rune, rows)
	for i := range b.grid {
		b.grid[i] = make([]rune, cols)
		for j := range b.grid[i] {
			b.grid[i][j] = '.'
		}
	}
	b.snakeCount = 0
	go b.generateFood()
}

func (b *Board) initSnake(s *Snake) {
	snakelength := 0
	i := b.rows / 2
	for j := 2; j < b.cols; j++ {
		s.body[snakelength] = Point{i + b.snakeCount, j}
		snakelength++
		if snakelength == MIN {
			s.head = Point{i + b.snakeCount, j}
			b.snakeCount++
			return
		}

	}
}

func (b *Board) update() {
	for _, s := range b.snakes {
		s.move()

		if s.head == b.food {
			s.score++
			// s.body = append([]Point{{}}, s.body...)
			b.generateFood()
		}
	}
}

func (b *Board) print() {
	var output strings.Builder

	for i := range b.rows {
		for j := range b.cols {
			printed := false
			if b.snakeCount != 0 {
				for _, snake := range b.snakes {
					for _, point := range snake.body {
						if point.x == i && point.y == j {
							output.WriteString("◉ ")
							printed = true
						}
					}
				}
			}
			if !printed {
				output.WriteString(string(b.grid[i][j]) + " ")
			}

		}
		output.WriteString("\n\r")
	}

	fmt.Print(output.String())
}

func clear() {
	fmt.Printf("\033[H\033[2J")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	conn     *websocket.Conn
	id       string
	keypress string
	snake    Snake
}

var (
	clients   = make(map[*Client]bool)
	clientsMu sync.Mutex
)

func (b *Board) run(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	client := &Client{
		conn: conn, id: r.RemoteAddr,
	}
	client.snake.init()

	clientsMu.Lock()
	b.snakes = append(b.snakes, &client.snake)
	b.initSnake(&client.snake)
	clientsMu.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Msg reading err:", err)
			break
		}

		client.keypress = string(msg)

		switch client.keypress {
		case "k":
			if client.snake.direction != down {
				client.snake.direction = up
			}
		case "j":
			if client.snake.direction != up {
				client.snake.direction = down
			}
		case "l":
			if client.snake.direction != left {
				client.snake.direction = right
			}

		case "h":
			if client.snake.direction != right {
				client.snake.direction = left
			}
		default:
			continue
		}

		fmt.Println("received:", client.keypress)

	}
}

func main() {
	b := &Board{}
	b.init(40, 60)
	clear()
	b.print()
	go func() {
		tick := time.NewTicker(tickRate)
		foodTick := time.NewTicker(foodPeriod)
		defer tick.Stop()
		defer foodTick.Stop()

		for {
			select {
			case <-tick.C:
				b.update()
				clear()
				b.print()
			}
		}
	}()

	http.HandleFunc("/", b.run)
	http.ListenAndServe(":8080", nil)
}
