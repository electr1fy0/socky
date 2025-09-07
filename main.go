package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"math/rand/v2"

	"golang.org/x/term"
)

type Direction int

// TODO:
// Implement food generation
// Implement boundary clashing
// Implement time based movement

const (
	up Direction = iota
	down
	left
	right
)

var over bool = false

var oldState *term.State

type Board struct {
	grid             [][]rune
	rows, cols       int
	snake            []Point
	tailRow, tailCol int
	headRow, headCol int
	snakeSize        int
	direction        Direction
	foodRow, foodCol int
}

type Point struct {
	x, y int
}

const MAX = 5

func (b *Board) init(rows, cols int) {
	b.rows, b.cols = rows, cols

	b.snake = make([]Point, MAX)
	snakeIndex := 0
	b.grid = make([][]rune, b.rows)
	for r := range b.grid {
		b.grid[r] = make([]rune, b.cols)
		for c := range b.grid[r] {
			if r == rows/2 && c > 2 && b.snakeSize != MAX {
				b.snake[snakeIndex] = Point{x: r, y: c}
				snakeIndex++
				b.snakeSize++
				if b.snakeSize == MAX {
					b.headRow = r
					b.headCol = c
				}
			}
			b.grid[r][c] = '.'
		}
	}

	x := rand.IntN(b.rows)
	y := rand.IntN(b.cols)

	b.foodRow = x
	b.foodCol = y

	b.grid[x][y] = '⊗'
	b.direction = right
}

func (b *Board) generateFood() {
	time.Sleep(4 * time.Second)
	x := rand.IntN(b.rows)
	y := rand.IntN(b.cols)

	b.foodRow = x
	b.foodCol = y

	b.grid[x][y] = '⊗'
}

func (b *Board) move() {
	switch b.direction {
	case up:
		b.moveUp()
	case down:
		b.moveDown()
	case left:
		b.moveLeft()
	case right:
		b.moveRight()
	}
}

func (b *Board) print() {
	var output strings.Builder
	for i := range b.rows {
		for j, val := range b.grid[i] {
			if i == b.headRow && j == b.headCol {
				output.WriteString("◕ ")
				continue
			}
			printed := false
			for _, point := range b.snake {
				if point.x == i && point.y == j {
					output.WriteString("◉ ")
					printed = true
					break
				}
			}
			if !printed {
				output.WriteString(string(val) + " ")
			}
		}
		output.WriteString("\n\r")
	}

	fmt.Print(output.String())
}

func (b *Board) moveDown() {

	if b.headCol >= b.cols || b.headRow >= b.rows || b.headRow <= 0 || b.headCol <= 0 {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Println("You lost")
		over = true
		return
	}

	b.headRow++
	n := len(b.snake)
	for i := 0; i < n-1; i++ {
		b.snake[i].x = b.snake[i+1].x
		b.snake[i].y = b.snake[i+1].y
	}
	b.snake[n-1].x = b.headRow
	if b.headRow == b.foodRow && b.headCol == b.foodCol {
		b.grid[b.foodRow][b.foodCol] = '.'
		b.snake = append([]Point{{x: b.tailRow, y: b.tailCol}}, b.snake...)
		go b.generateFood()
	}
}

func (b *Board) moveUp() {

	if b.headCol >= b.cols || b.headRow >= b.rows || b.headRow <= 0 || b.headCol <= 0 {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Println("You lost")
		over = true
		return
	}

	b.headRow--
	n := len(b.snake)
	b.tailCol = b.snake[0].y
	b.tailRow = b.snake[0].x

	for i := 0; i < n-1; i++ {
		b.snake[i].x = b.snake[i+1].x
		b.snake[i].y = b.snake[i+1].y
	}
	b.snake[n-1].x = b.headRow

	if b.headRow == b.foodRow && b.headCol == b.foodCol {
		b.grid[b.foodRow][b.foodCol] = '.'
		b.snake = append([]Point{{x: b.tailRow, y: b.tailCol}}, b.snake...)
		go b.generateFood()
	}
}

func (b *Board) moveLeft() {
	if b.headCol >= b.cols || b.headRow >= b.rows || b.headRow <= 0 || b.headCol <= 0 {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Println("You lost")
		over = true
		return
	}

	b.headCol--
	n := len(b.snake)
	b.tailCol = b.snake[0].y
	b.tailRow = b.snake[0].x

	for i := 0; i < n-1; i++ {
		b.snake[i].x = b.snake[i+1].x
		b.snake[i].y = b.snake[i+1].y
	}
	b.snake[n-1].y = b.headCol

	if b.headRow == b.foodRow && b.headCol == b.foodCol {
		b.grid[b.foodRow][b.foodCol] = '.'
		b.snake = append([]Point{{x: b.tailRow, y: b.tailCol}}, b.snake...)

		go b.generateFood()
	}
}

func (b *Board) moveRight() {
	if b.headCol >= b.cols || b.headRow >= b.rows || b.headRow <= 0 || b.headCol <= 0 {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Println("You lost")
		over = true
		return
	}

	b.headCol++
	n := len(b.snake)
	b.tailCol = b.snake[0].y
	b.tailRow = b.snake[0].x
	for i := 0; i < n-1; i++ {
		b.snake[i].x = b.snake[i+1].x
		b.snake[i].y = b.snake[i+1].y
	}
	b.snake[n-1].y = b.headCol
	if b.headRow == b.foodRow && b.headCol == b.foodCol {
		b.grid[b.foodRow][b.foodCol] = '.'
		b.snake = append([]Point{{x: b.tailRow, y: b.tailCol}}, b.snake...)
		go b.generateFood()

	}
}

func clear() {
	fmt.Printf("\033[H\033[2J")
}

func main() {
	b := Board{}
	b.init(40, 60)
	clear()
	b.print()
	oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	// if err != nil {
	// panic(err)
	// }

	buf := make([]byte, 1)
	var move string
	go func() {
		for {
			os.Stdin.Read(buf)
			move = strings.ToLower(string(buf))
			switch move {
			case "k":
				if b.direction != down {
					b.direction = up
				}
			case "j":
				if b.direction != up {
					b.direction = down
				}
			case "l":
				if b.direction != left {
					b.direction = right
				}

			case "h":
				if b.direction != right {
					b.direction = left
				}
			default:
				over = true
				return
			}
		}
	}()

	for {
		if over {
			break
		}
		b.move()
		clear()
		b.print()
		time.Sleep(100 * time.Millisecond)
	}
}
