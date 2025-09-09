package main

// This is not a part of the server module
// Run this independently to connect to the servers

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"golang.org/x/term"
)

var oldState *term.State

var err error
var url = "ws://localhost:8081"
var conn *websocket.Conn

func clear() {
	fmt.Printf("\033[H\033[2J")
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Error dialing up:", err)
	}
	defer func() {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing from client"))
		conn.Close()
	}()
	conn.WriteMessage(websocket.TextMessage, []byte("client has connected"))

	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1)
	over := false

	go func() {
		for {
			os.Stdin.Read(buf)
			if string(buf) == "q" {
				over = true
				return
			}
			conn.WriteMessage(websocket.TextMessage, buf)
		}
	}()

	for {
		if over {
			return
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading msg:", err)
		}
		clear()
		fmt.Print(string(msg))
	}

}
