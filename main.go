package main

import (
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
	// var move string
	// go func() {
	for {
		os.Stdin.Read(buf)
		// move = strings.ToLower(string(buf))
		conn.WriteMessage(websocket.TextMessage, buf)
		// _, msg, _ := conn.ReadMessage()
		// fmt.Println(string(msg))
		if string(buf) == "q" {
			return
		}
	}
	// }()

	// for {
	// 	if over {
	// 		if lost {
	// 			fmt.Println("You lose")
	// 		}
	// 		fmt.Println("Score: ", b.score)
	// 		break
	// 	}
	// 	b.move()
	// 	clear()
	// 	b.print()
	// 	time.Sleep(200 * time.Millisecond)
	// }
}
