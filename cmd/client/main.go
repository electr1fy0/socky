package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/coder/websocket"
	"golang.org/x/term"
)

var (
	oldState *term.State
	err      error
	url      = "ws://localhost:8080/"
	name     string
	game     Message
)

const (
	writeWait = time.Second * 5
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "internet" {
		if u := os.Getenv("SOCKY_SERVER_URL"); u != "" {
			url = u
		}
	}
	Clear()

	WelcomeBanner()
	fmt.Print("\t Enter your name (keep it short): ")
	fmt.Scanln(&name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	conn, _, err := websocket.Dial(ctx, url, nil)
	defer cancel()
	if err != nil {
		fmt.Println("Server is not ready.")
		os.Exit(1)
	}
	defer func() {
		conn.Close(websocket.StatusNormalClosure, "closing from client")
	}()

	if err = conn.Write(ctx, websocket.MessageText, []byte("client has connected")); err != nil {
		fmt.Println("Error writing message:", err)
		os.Exit(1)
	}
	if err = conn.Write(ctx, websocket.MessageText, []byte("NAME:"+name)); err != nil {
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
				conn.Close(websocket.StatusNormalClosure, "closing from client")
				return
			}
			conn.Write(ctx, websocket.MessageText, buf)
		}
	}()

	for {
		_, msg, err := conn.Read(ctx)

		if err != nil {
			GameOverBanner()
			os.Exit(0)
		}
		json.Unmarshal(msg, &game)
		if game.Type == "over" {
			GameOverBanner()
			os.Exit(0)
		}
		Clear()
		fmt.Print(RenderGame())
	}
}
