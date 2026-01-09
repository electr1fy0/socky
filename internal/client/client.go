package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/electr1fy0/socky/internal/game"
	"golang.org/x/term"
)

var (
	err           error
	wsBaseURL     = "ws://localhost:8080"
	httpBaseURL   = "http://localhost:8080"
	stateRestored = false

	name    string
	message game.Message
)

const (
	writeWait = time.Second * 5
)

var ctx, cancel = context.WithTimeout(context.Background(), time.Minute)

func sendKeystrokes(conn *websocket.Conn) {
	buf := make([]byte, 1)

	for {
		os.Stdin.Read(buf)
		switch string(buf) {
		case "q":
			conn.Close(websocket.StatusNormalClosure, "closing from client")
			return
		case "c":
			CreateRoom()
		}
		conn.Write(ctx, websocket.MessageText, buf)
	}
}

func CreateRoom() string {
	resp, err := http.Get(httpBaseURL + "/create")
	if err != nil {
		fmt.Println("failed to create a room:", err)
		return ""
	}
	defer resp.Body.Close()
	var data struct{ ID string }

	json.NewDecoder(resp.Body).Decode(&data)
	return data.ID

}

func JoinRoom(id string) (*websocket.Conn, error) {
	conn, _, err := websocket.Dial(ctx, wsBaseURL+"/join/"+id, nil)

	return conn, err
}

func InitiateGame() {
	id := os.Args[1]
	if "create" == strings.TrimSpace(id) {
		newRoomID := CreateRoom()
		fmt.Println("Room ID:", newRoomID)
		fmt.Println("Share this ID with a friend to join the game.")
		return
	}

	Clear()

	WelcomeBanner()
	fmt.Print("\t Enter your name (keep it short): ")
	fmt.Scanln(&name)

	defer cancel()
	if err != nil {
		fmt.Println("Server is not ready.")
		os.Exit(1)
	}

	conn, err := JoinRoom(id)
	if err != nil {
		fmt.Println("Failed to join room:", err)
		os.Exit(1)
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing from client")

	if err := conn.Write(ctx, websocket.MessageText, []byte("client has connected")); err != nil {
		fmt.Println("Error writing message:", err)
		os.Exit(1)
	}

	if err := conn.Write(ctx, websocket.MessageText, []byte("NAME:"+name)); err != nil {
		fmt.Println("Error sending name:", err)
		os.Exit(1)
	}

	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	defer func() {
		if !stateRestored {
			term.Restore(int(os.Stdin.Fd()), oldState)
		}
	}()
	if err != nil {
		panic(err)
	}
	go sendKeystrokes(conn)

	for {
		_, msg, err := conn.Read(ctx)

		if err != nil {
			GameOverBanner()
			return
		}
		json.Unmarshal(msg, &message)
		if message.Type == "over" {
			GameOverBanner()
			return
		}
		Clear()
		fmt.Print(RenderGame())
	}
}
