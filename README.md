# Socky

Terminal-based multiplayer Snake game built with Go and WebSockets.

## Screenshots
<div style="gap:10px;">
<img style="width: 60%; height="auto" alt="Screenshot 2025-09-18 at 11 34 50" src="https://github.com/user-attachments/assets/defde25f-ddb0-4ff0-a5e4-ba9ef4664f76" />
<img style="width: 60%; height="auto" alt="Screenshot 2025-09-18 at 11 41 03" src="https://github.com/user-attachments/assets/e89495ae-e4e8-4ef9-9828-6f398c3f6db6" />
</div>


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
