package game

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
)

type Point struct{ X, Y int }

type Board struct {
	Rows, Cols  int
	Grid        [][]rune
	SnakeCount  int
	Clients     []*Client
	clientCount int
	mu          sync.RWMutex
}

func (b *Board) GenerateFood() {
	x := rand.IntN(b.Rows)
	y := rand.IntN(b.Cols)

	b.mu.Lock()
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
				b.Grid[i][j] = '·'
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

		prevHead := c.Snake.Body[len(c.Snake.Body)-2]
		b.Grid[prevHead.X][prevHead.Y] = '◉'
		b.Grid[c.Snake.Tail.X][c.Snake.Tail.Y] = '·'

		if b.Grid[c.Snake.Head.X][c.Snake.Head.Y] == '◆' {
			c.Snake.Score++

			c.Snake.Body = append([]Point{c.Snake.Tail}, c.Snake.Body...)
		}
		b.Grid[c.Snake.Head.X][c.Snake.Head.Y] = '◕'

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
	output.WriteString("┐\t\r\n")

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
		output.WriteString("│\t\r\n")
	}
	output.WriteString("\t")
	output.WriteString("└")
	for i := 0; i <= b.Cols*2; i++ {
		output.WriteString("─")
	}
	output.WriteString("┘\t\r\n")

	return output.String()
}

func Clear() {
	fmt.Printf("\033[H\033[2J")
}
