package main

import (
	"net/http"
	"time"

	"github.com/electr1fy0/socky/board"
)

// TODO:
// 1. Write board back to the clients
// 2. Leaderboard
// 3. Broadcast fn
// 4. Graceful connection closure and snake removal upn connection(!!!)
//

func main() {
	b := &board.Board{}
	b.Init(40, 60)
	b.Print()
	go func() {
		tick := time.NewTicker(board.TickRate)
		foodTick := time.NewTicker(board.FoodPeriod)
		defer tick.Stop()
		defer foodTick.Stop()

		for { // reminder: modernize this after understanding range in channel
			select {
			case <-tick.C:
				b.Update()
			}
		}
	}()

	http.HandleFunc("/", b.Run)
	http.ListenAndServe(":8080", nil)
}
