package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/coder/websocket"

	"github.com/electr1fy0/socky/internal/game"
)

func getKeypresses(client *game.Client, ctx context.Context) {
	for {
		_, msg, err := client.Conn.Read(ctx)

		if err != nil {
			fmt.Println("Err reading:", err)
			break
		}
		msgString := string(msg)
		if name, ok := strings.CutPrefix(msgString, "NAME:"); ok {
			client.Name = name
			continue
		}

		client.Keypress = string(msg)

		switch client.Keypress {
		case "k", "w":
			if client.Snake.Direction != game.Down {
				client.Snake.Direction = game.Up
			}
		case "j", "s":
			if client.Snake.Direction != game.Up {
				client.Snake.Direction = game.Down
			}
		case "l", "d":
			if client.Snake.Direction != game.Left {
				client.Snake.Direction = game.Right
			}
		case "h", "a":
			if client.Snake.Direction != game.Right {
				client.Snake.Direction = game.Left
			}
		default:
			continue
		}
		fmt.Print("received:", client.Keypress, "\r")
	}
}

func (R *RoomManager) HandleCreate(w http.ResponseWriter, r *http.Request) {
	room := NewRoom()

	R.Mu.Lock()
	R.Rooms[room.ID] = room
	R.Mu.Unlock()
	go room.InitGame()
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"ID": room.ID,
	}
	json.NewEncoder(w).Encode(resp)
}

func (R *RoomManager) HandleJoin(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomID")
	fmt.Println("got id:", roomID)
	room := R.Rooms[roomID]
	if room == nil {
		fmt.Println("nil room")
		return
	}
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Println("Error: ", err)
		return
	}

	client := &game.Client{
		Conn: conn, ID: r.RemoteAddr,
	}

	defer func() {
		room.Board.RemoveClient(client)
		conn.Close(websocket.StatusNormalClosure, "client left")
	}()

	client.Snake.Init()
	room.Board.AddClient(client)

	getKeypresses(client, context.Background())
}
