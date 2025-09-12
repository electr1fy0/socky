package main

// This is not a part of the server module
// Run this independently to connect to the servers

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"golang.org/x/term"
)

var oldState *term.State

var err error
var url = "ws://localhost:8080/"

var conn *websocket.Conn

func clear() {
	fmt.Printf("\033[H\033[2J")
}

func main() {
	// if u := os.Getenv("SOCKY_SERVER_URL"); u != "" {
	// 	url = u
	// }
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println("Error dialing up:", err)
		os.Exit(1)
	}

	var name string
	fmt.Print("Enter your name: ")
	fmt.Scanln(&name)

	defer func() {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing from client"))
		conn.Close()
	}()
	err = conn.WriteMessage(websocket.TextMessage, []byte("client has connected"))
	if err != nil {
		fmt.Println("Error writing message:", err)
		os.Exit(1)
	}

	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	if err != nil {
		panic(err)
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("NAME:"+name))
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
			// reminder: come on, mate. you're better than this (i'm ashamed)
			fmt.Println("\n\r\tYOU LOST!")
			break
		}
		clear()
		fmt.Print(string(msg))
		fmt.Print("\t<hjkl> or <wasd> to move. <q> to quit.\n\r")
	}

}
