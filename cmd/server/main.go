package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/electr1fy0/socky/game.go"
)

// TODO:
// 1. Write board back to the clients
// 2. Leaderboard
// 3. Broadcast fn
// 4. Graceful connection closure and snake removal upn connection(!!!)
// 5. Name the players
// 6. Fix double close on clash with wall, unify closing logic
// 7. Use JSON, probably

const (
	boardHeight = 30
	boardWidth  = 35
	FoodPeriod  = 8 * time.Second
	TickRate    = 200 * time.Millisecond
)

func main() {
	b := &game.Board{}
	b.Init(boardHeight, boardWidth)
	b.Print()
	go func() {
		tick := time.NewTicker(TickRate)
		foodTick := time.NewTicker(FoodPeriod)
		defer tick.Stop()
		defer foodTick.Stop()
		for {
			select {
			case <-tick.C:
				b.Update()
				b.ClearBoard()
				b.BroadCast()
			case <-foodTick.C:
				b.GenerateFood()
			}
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
