package main

import (
	"fmt"
	"net/http"
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

	http.HandleFunc("/", b.Run)
	fmt.Println("Server is up at port 8081")

	http.ListenAndServe(":8081", nil)
}
