package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	boardHeight = 30
	boardWidth  = 35
	FoodPeriod  = 8 * time.Second
	TickRate    = 150 * time.Millisecond
)

type Room struct {
	ID string

	Board  *Board
	ticker *time.Ticker
	done   chan struct{}
}

type RoomManager struct {
	Mu    sync.RWMutex
	Rooms map[string]*Room
}

func NewRoom() *Room {
	b := &Board{}
	b.Init(boardHeight, boardWidth)
	id := uuid.NewString()
	fmt.Println("Room ID: ", id)

	return &Room{
		ID:    id,
		Board: b,
		done:  make(chan struct{}),
	}
}

func (R *Room) InitRoom() {
	tick := time.Tick(TickRate)
	foodTick := time.Tick(FoodPeriod)
	for {
		select {
		case <-tick:
			R.Board.Update()
			R.Board.ClearBoard()
			R.Board.BroadCast()
		case <-foodTick:
			R.Board.GenerateFood()
		case <-R.done:
			return
		}
	}
}

func (R *Room) Close() {
	select {
	case <-R.done:
	default:
		close(R.done)
	}
}
