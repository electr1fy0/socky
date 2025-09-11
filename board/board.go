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
	MinSnakeLength = 5
	FoodPeriod     = 8 * time.Second
	TickRate       = 200 * time.Millisecond
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

var snakeColors = []string{
	"\033[31m", // Red
	"\033[33m", // Yellow
	"\033[32m", // Green
	"\033[34m", // Blue
	"\033[35m", // Magenta
	"\033[36m", // Cyan
	"\033[37m", // White
}

var foodColor = "\033[38;5;208m" // Bright orange
const resetColor = "\033[0m"

type Point struct{ X, Y int }

type Snake struct {
	Body      []Point
	Head      Point
	Tail      Point
	Direction Direction
	Score     int
}

func (s *Snake) Init() {
	s.Body = make([]Point, MinSnakeLength)
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
	Grid       [][]rune
	SnakeCount int
	Clients    []*Client
	mu         sync.RWMutex
}

func (b *Board) GenerateFood() {
	x := rand.IntN(b.Rows)
	y := rand.IntN(b.Cols)

	b.mu.Lock()
	b.Food = Point{X: x, Y: y}
	b.Grid[x][y] = '◆'
	b.mu.Unlock()
}

func (b *Board) Init(rows, cols int) {
	b.Rows, b.Cols = rows, cols
	b.Clients = []*Client{}
	b.Grid = make([][]rune, rows)
	for i := range b.Grid {
		b.Grid[i] = make([]rune, cols)
		for j := range b.Grid[i] {
			b.Grid[i][j] = '·'
		}
	}
	b.SnakeCount = 0
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
		if snakelength == MinSnakeLength {
			s.Head = Point{i + b.SnakeCount, j}
			b.SnakeCount++
			return
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
			for _, otherBody := range other.Snake.Body {
				if c.Snake.Head == otherBody {
					toRemove = append(toRemove, c)
					added = true
					break
				}
			}
		}
		if added {
			continue
		}
		if c.Snake.Head.X < 0 || c.Snake.Head.X >= b.Rows || c.Snake.Head.Y < 0 || c.Snake.Head.Y >= b.Cols {
			toRemove = append(toRemove, c)
			continue
		}

		b.Grid[c.Snake.Head.X][c.Snake.Head.Y] = '◕'

		prevHead := c.Snake.Body[len(c.Snake.Body)-2]
		b.Grid[prevHead.X][prevHead.Y] = '◉'
		b.Grid[c.Snake.Tail.X][c.Snake.Tail.Y] = '·'

		if c.Snake.Head == b.Food {
			c.Snake.Score++
			c.Snake.Body = append([]Point{c.Snake.Tail}, c.Snake.Body...)
			go b.GenerateFood()
		}
	}
	for _, c := range toRemove {
		b.removeClient(c)
		c.Conn.Close()
	}
}

func (b *Board) Print() string {
	var output strings.Builder
	output.WriteString("\t┌")
	for i := 0; i <= b.Cols*2; i++ {
		output.WriteString("─")
	}
	output.WriteString("┐\r\n")

	for i := 0; i < b.Rows; i++ {
		output.WriteString("\t│ ")
		for j := 0; j < b.Cols; j++ {
			cellSymbol := string(b.Grid[i][j])
			colored := cellSymbol

			if b.Grid[i][j] == '◆' {
				colored = foodColor + string('◆') + resetColor
				output.WriteString(colored + " ")
				continue
			}
			for _, client := range b.Clients {
				for _, body := range client.Snake.Body {
					if body.X == i && body.Y == j {
						colored = client.Color + cellSymbol + resetColor
						break
					}
				}
			}
			output.WriteString(colored + " ")
		}
		output.WriteString("│\r\n")
	}
	output.WriteString("\t")
	output.WriteString("└")
	for i := 0; i <= b.Cols*2; i++ {
		output.WriteString("─")
	}
	output.WriteString("┘\r\n")

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
	Name     string
	Color    string
}

func (b *Board) addClient(client *Client) {
	client.Color = snakeColors[len(b.Clients)%len(snakeColors)]
	b.mu.Lock()
	b.Clients = append(b.Clients, client)
	b.mu.Unlock()
}

func (b *Board) BroadCast() {
	b.mu.RLock()
	boardState := b.Print()
	clients := make([]*Client, len(b.Clients))

	copy(clients, b.Clients)

	b.mu.RUnlock()

	for _, client := range clients {
		scoreText := "\tScores:\r"
		for _, currentClient := range clients {
			name := client.Name
			scoreText += "\n\r\t" + name + ": " + strconv.Itoa(currentClient.Snake.Score) + "\r\n"
		}

		if err := client.Conn.WriteMessage(websocket.TextMessage, []byte(boardState+scoreText+"\n\n")); err != nil {
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
					b.Grid[point.X][point.Y] = '○'
				}
			}
			b.Grid[c.Snake.Tail.X][c.Snake.Tail.Y] = '○'
			b.Clients = append(b.Clients[:i], b.Clients[i+1:]...)
			break
		}
	}
	b.mu.Unlock()
	go func() {
		time.Sleep(1 * time.Second)
		b.mu.Lock()
		for i := range b.Grid {
			for j := range b.Grid[i] {
				if b.Grid[i][j] == '○' {
					b.Grid[i][j] = '·'
				}
			}
		}
		b.mu.Unlock()
	}()
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
	b.addClient(client)
	b.InitSnake(&client.Snake)

	defer func() {
		if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "see ya, mate")); err != nil {
			fmt.Println("Error writing close message: (ignore for now)", err)
		}
		b.removeClient(client)

		fmt.Println("client left")
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()

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
