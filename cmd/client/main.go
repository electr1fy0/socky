package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/term"
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
	foodColor  = "\033[38;5;208m" // Bright orange
	resetColor = "\033[0m"
	oldState   *term.State
	err        error
	url        = "ws://localhost:8080/"
	name       string
	game       Message
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	Body  []Point
	Head  Point
	Tail  Point
	Score int
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
	Grid    [][]string `json:"grid"`
	Clients []*Client  `json:"clients"`
}

func clear() {
	fmt.Printf("\033[H\033[2J")
}

func Print() string {
	var output strings.Builder

	output.WriteString("\t┌")
	for i := 0; i <= len(game.Grid[0])*2; i++ {
		output.WriteString("─")
	}
	output.WriteString("┐\t\r\n")

	for i := 0; i < len(game.Grid); i++ {
		output.WriteString("\t│ ")
		for j := 0; j < len(game.Grid[0]); j++ {
			cellSymbol := string(game.Grid[i][j])
			colored := cellSymbol

			if cellSymbol == "f" {
				output.WriteString(foodColor + "f " + resetColor)
				continue
			}

			for _, client := range game.Clients {
				for _, body := range client.Snake.Body {
					if body.X == i && body.Y == j {
						colored = snakeColors[client.Color] + cellSymbol + resetColor
						break
					}
				}
			}

			output.WriteString(colored + " ")
		}
		output.WriteString("│\t\r\n")
	}

	output.WriteString("\t└")
	for i := 0; i <= len(game.Grid[0])*2; i++ {
		output.WriteString("─")
	}
	output.WriteString("┘\t\r\n")

	return output.String()
}

func getScore() string {
	bold := "\033[1m"
	scoreText := "\n\t " + bold + "------------------------------- SCORES --------------------------------" + resetColor + "\r\n"

	for _, client := range game.Clients {
		bar := strings.Repeat("█", client.Snake.Score)
		bar = snakeColors[client.Color] + bar + resetColor
		scoreText += fmt.Sprintf("\t %-8s | %3d  %s\n", client.Name, client.Snake.Score, bar)
	}

	scoreText += "\n"
	return scoreText
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "internet" {
		if u := os.Getenv("SOCKY_SERVER_URL"); u != "" {
			url = u
		}
	}

	fmt.Print("Enter your name (keep it short): ")
	fmt.Scanln(&name)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println("Error dialing up:", err)
		os.Exit(1)
	}
	defer func() {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing from client"))
		conn.Close()
	}()

	if err = conn.WriteMessage(websocket.TextMessage, []byte("client has connected")); err != nil {
		fmt.Println("Error writing message:", err)
		os.Exit(1)
	}
	if err = conn.WriteMessage(websocket.TextMessage, []byte("NAME:"+name)); err != nil {
		fmt.Println("Error sending name:", err)
		os.Exit(1)
	}

	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1)
	go func() {
		for {
			os.Stdin.Read(buf)
			if string(buf) == "q" {
				conn.Close()
				return
			}
			conn.WriteMessage(websocket.TextMessage, buf)
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("\n\r\tYou lost!")
			break
		}
		json.Unmarshal(msg, &game)
		clear()
		fmt.Print(Print())
		fmt.Print(getScore())
		fmt.Print("\r\t<hjkl> or <wasd> to move. <q> to quit.\n\r")
	}
}
