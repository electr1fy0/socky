package main

import "github.com/coder/websocket"

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	Body  []Point
	Head  Point
	Tail  Point
	Score int
}

type Client struct {
	ID       string          `json:"id"`
	Keypress string          `json:"keypress"`
	Snake    Snake           `json:"snake"`
	Name     string          `json:"name"`
	Color    string          `json:"color"`
	Conn     *websocket.Conn `json:"-"`
}

type Message struct {
	Type    string     `json:"type"`
	Grid    [][]string `json:"grid"`
	Clients []*Client  `json:"clients"`
}
