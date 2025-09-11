# Socky

Terminal-based multiplayer Snake game built with Go and WebSockets.

## Features
- **Real-time multiplayer** - Multiple players compete on the same board
- **Terminal graphics** - Runs entirely in your terminal with Unicode symbols
- **Live collision detection** - Players die when hitting walls or each other
- **Live scoreboard** - See all players' scores update in real-time

## Usage
```bash
# Start server
go run main.go

# Connect players (run multiple times for multiplayer)
go run client.go
```

## Controls
- **WASD** or **HJKL** (vim keys) to move

## Dependencies
- `github.com/gorilla/websocket`
