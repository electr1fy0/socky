package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func Run() {
	roomManager := &RoomManager{
		Rooms: make(map[string]*Room),
	}

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/create", roomManager.HandleCreate)
	http.HandleFunc("GET /join/{roomID}", roomManager.HandleJoin)

	fmt.Printf("Starting server at port %s...\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
