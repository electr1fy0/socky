package board

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	Min        = 5
	TickRate   = 600 * time.Millisecond
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
	ID        string
	haslost   bool
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
	Clients    []*Client
	mu         sync.RWMutex
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
			go b.GenerateFood()
		}
	}
}

func (b *Board) Print() string {
	var output strings.Builder

	for i := 0; i < b.Rows; i++ {
		for j := 0; j < b.Cols; j++ {
			printed := false
			if b.SnakeCount != 0 {
				for _, snake := range b.Snakes {
					if i == snake.Head.X && j == snake.Head.Y {
						output.WriteString("◕ ")
						printed = true
						continue
					}
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

	return output.String()
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

func (b *Board) addClient(client *Client) {
	b.mu.Lock()
	b.Clients = append(b.Clients, client)
	b.Snakes = append(b.Snakes, &client.Snake)
	b.mu.Unlock()
}

func (b *Board) BroadCast() {
	b.mu.RLock()
	boardState := b.Print()
	snakes := make([]*Snake, len(b.Snakes))
	copy(snakes, b.Snakes)
	clients := make([]*Client, len(b.Clients))

	copy(clients, b.Clients)

	b.mu.RUnlock()

	for _, client := range clients {
		scoreText := "Scores:\r"
		otherCnt := 0
		for _, snake := range snakes {
			if snake.ID == client.Snake.ID {
				scoreText += "\n\rYou: " + strconv.Itoa(snake.Score) + "\r"
			} else {
				otherCnt++
				scoreText += "\n\rPlayer " + strconv.Itoa(otherCnt) + ": " + strconv.Itoa(snake.Score) + "\r"
			}
		}
		client.Conn.WriteMessage(websocket.TextMessage, []byte(boardState+scoreText))
	}
}

func (b *Board) removeClient(client *Client) {
	b.mu.Lock()
	for i, s := range b.Snakes {
		if s.ID == client.Snake.ID {
			b.Snakes = append(b.Snakes[:i], b.Snakes[i+1:]...)
			break
		}
	}
	for i, c := range b.Clients {
		if c.Snake.ID == client.Snake.ID {
			b.Clients = append(b.Clients[:i], b.Clients[i+1:]...)
			break
		}
	}
	b.mu.Unlock()
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

	client.Snake.Init()
	client.Snake.ID = r.RemoteAddr
	b.addClient(client)
	b.InitSnake(&client.Snake)

	defer func() {
		b.removeClient(client)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "see ya, mate"))
		fmt.Println("client left")
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Err reading:", err)
			break
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
