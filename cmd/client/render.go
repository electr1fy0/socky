package main

import (
	"fmt"
	"slices"
	"strings"
)

var (
	snakeColors = map[string]string{
		"red":     "\033[31m",
		"yellow":  "\033[33m",
		"green":   "\033[32m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"white":   "\033[37m",
	}
	foodColor  = "\033[38;5;208m"
	resetColor = "\033[0m"
)

func Clear() {
	fmt.Printf("\033[H\033[2J")
}

func RenderGame() string {
	var output strings.Builder

	output.WriteString("\n\n\t╔")
	for i := 0; i <= len(game.Grid[0])*2; i++ {
		output.WriteString("═")
	}
	output.WriteString("╗\t\r\n")

	shadowColor := "\033[38;5;240m"
	board := make([][]string, len(game.Grid))
	for i := range board {
		board[i] = make([]string, len(game.Grid[0]))
		for j := range board[i] {
			board[i][j] = "  "
		}
	}

	for i := 0; i < len(game.Grid); i++ {
		for j := 0; j < len(game.Grid[0]); j++ {
			if game.Grid[i][j] == "f" {
				board[i][j] = foodColor + "◆ " + resetColor
			}
		}
	}

	for _, client := range game.Clients {
		color := snakeColors[client.Color]
		for _, body := range client.Snake.Body {
			symbol := "█ "
			if body == client.Snake.Head {
				symbol = "◉ "
			}

			board[body.X][body.Y] = color + symbol + resetColor

			if body.X+1 < len(board) && body.Y+1 < len(board[0]) {
				if board[body.X+1][body.Y+1] == "  " {
					board[body.X+1][body.Y+1] = shadowColor + "░ " + resetColor
				}
			}
		}
	}

	for i := range len(board) {
		output.WriteString("\t║ ")
		for j := 0; j < len(board[0]); j++ {
			output.WriteString(board[i][j])
		}
		output.WriteString("║\t\r\n")
	}

	output.WriteString("\t╚")
	for i := 0; i <= len(board[0])*2; i++ {
		output.WriteString("═")
	}
	output.WriteString("╝\t\r\n" + renderScoreboard() + "\r\t<hjkl> or <wasd> to move. <q> to quit.\n\r")

	return output.String()
}

func renderScoreboard() string {
	scoreBuilder := strings.Builder{}
	bold := "\033[1m"
	scoreBuilder.WriteString("\n\t " + bold + "------------------------------ SCOREBOARD ------------------------------" + resetColor + "\r\n\n")

	clients := make([]*Client, len(game.Clients))
	copy(clients, game.Clients)

	slices.SortFunc(clients, func(a, b *Client) int {
		return b.Snake.Score - a.Snake.Score
	})

	for rank, client := range clients {
		bar := strings.Repeat("█", client.Snake.Score)
		bar = snakeColors[client.Color] + bar + resetColor
		crown := ""
		if rank == 0 {
			crown = " 👑"
		}

		fmt.Fprintf(&scoreBuilder, "\t %2d. %-8s | %3d  %s%s\r\n", rank+1, client.Name, client.Snake.Score, bar, crown)
	}

	scoreBuilder.WriteString("\n")
	return scoreBuilder.String()
}
