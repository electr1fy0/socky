package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type Direction int

const (
	top Direction = iota
	down
	right
	left
)

type Board struct {
	grid             [][]rune
	rows, cols       int
	snake            []Point
	head, tail, body int
	headRow, headCol int
	snakeSize        int
}

type Point struct {
	x, y int
}

const MAX = 25

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
}

func (b Board) print() {
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
	b.headRow++
	n := len(b.snake)
	for i := 0; i < n-1; i++ {
		b.snake[i].x = b.snake[i+1].x
		b.snake[i].y = b.snake[i+1].y
	}
	b.snake[n-1].x = b.headRow
}

func (b *Board) moveUp() {
	b.headRow--
	n := len(b.snake)
	for i := 0; i < n-1; i++ {
		b.snake[i].x = b.snake[i+1].x
		b.snake[i].y = b.snake[i+1].y
	}
	b.snake[n-1].x = b.headRow
}

func (b *Board) moveLeft() {
	b.headCol--
	n := len(b.snake)
	for i := 0; i < n-1; i++ {
		b.snake[i].x = b.snake[i+1].x
		b.snake[i].y = b.snake[i+1].y
	}
	b.snake[n-1].y = b.headCol
}

func (b *Board) moveRight() {
	b.headCol++
	n := len(b.snake)
	for i := 0; i < n-1; i++ {
		b.snake[i].x = b.snake[i+1].x
		b.snake[i].y = b.snake[i+1].y
	}
	b.snake[n-1].y = b.headCol
}

func clear() {
	fmt.Printf("\033[H\033[2J")
}

func main() {
	b := Board{}
	b.init(40, 60)
	clear()
	b.print()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 1)
	var move string

	for {
		os.Stdin.Read(buf)
		move = strings.ToLower(string(buf[0]))
		switch move {
		case "k":
			b.moveUp()
		case "j":
			b.moveDown()
		case "l":
			b.moveRight()
		case "h":
			b.moveLeft()
		case "q":
			return
		default:
			fmt.Println("invalid move")
			os.Exit(1)
		}
		clear()
		b.print()
	}
}
