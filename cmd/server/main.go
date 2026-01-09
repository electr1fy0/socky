package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/electr1fy0/socky/internal"
)

const (
	boardHeight = 30
	boardWidth  = 35
	FoodPeriod  = 8 * time.Second
	TickRate    = 150 * time.Millisecond
)

func main() {

	room := internal.NewRoom()
	go room.InitRoom()

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/", room.Board.Run)
	fmt.Printf("Starting server at port %s...\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
