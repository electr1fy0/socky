package main

import (
	"net/http"
	"time"

	"github.com/electr1fy0/socky/board"
)

func main() {
	b := &board.Board{}
	b.Init(40, 60)
	b.Print()
	go func() {
		tick := time.NewTicker(board.TickRate)
		foodTick := time.NewTicker(board.FoodPeriod)
		defer tick.Stop()
		defer foodTick.Stop()

		for {
			select {
			case <-tick.C:
				b.Update()
				board.Clear()
				b.Print()
			}
		}
	}()

	http.HandleFunc("/", b.Run)
	http.ListenAndServe(":8080", nil)
}
