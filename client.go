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

var over, lost bool = false, false

var oldState *term.State

func connect() {

}

var err error
var url = "ws://localhost:8080"
var conn *websocket.Conn

func clear() {
	fmt.Printf("\033[H\033[2J")
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Error dialing up:", err)
	}

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
			conn.WriteMessage(websocket.TextMessage, buf)
			if string(buf) == "q" {
				over = true
				return
			}

		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading msg:", err)
		}
		clear()
		fmt.Print(string(msg))
		if over {
			break
		}
	}

}
