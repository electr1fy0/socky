package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/electr1fy0/socky/game"
)

// TODO:
// 1. Fix double close on clash with wall, unify closing logic

const (
	boardHeight = 30
	boardWidth  = 35
	FoodPeriod  = 8 * time.Second
	TickRate    = 150 * time.Millisecond
)

func main() {
	b := &game.Board{}
	b.Init(boardHeight, boardWidth)

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
	fmt.Printf("Starting server at port %s...\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
