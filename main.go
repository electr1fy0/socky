package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/electr1fy0/socky/board"
)

// TODO:
// 1. Write board back to the clients
// 2. Leaderboard
// 3. Broadcast fn
// 4. Graceful connection closure and snake removal upn connection(!!!)
// 5. Name the players

const boardHeight = 40
const boardWidth = 50

func main() {
	b := &board.Board{}
	b.Init(boardHeight, boardWidth)
	b.Print()
	go func() {
		tick := time.NewTicker(board.TickRate)
		foodTick := time.NewTicker(board.FoodPeriod)
		defer tick.Stop()
		defer foodTick.Stop()

		for range tick.C {
			b.Update()
			b.BroadCast()
		}
	}()

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/", b.Run)
	fmt.Println("Server is up at port", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error listening: ", err)
	}
}
