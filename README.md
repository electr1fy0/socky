# Socky

Terminal-based multiplayer Snake game built with Go and WebSockets.

## Screenshots
<div style="gap:10px;">
<img style="width: 60%; height="auto" alt="Screenshot 2025-09-18 at 11 34 50" src="https://github.com/user-attachments/assets/defde25f-ddb0-4ff0-a5e4-ba9ef4664f76" />
<img style="width: 60%; height="auto" alt="Screenshot 2025-09-18 at 11 41 03" src="https://github.com/user-attachments/assets/e89495ae-e4e8-4ef9-9828-6f398c3f6db6" />
</div>


## Features
- Real-time multiplayer: multiple players on the same board
- Terminal graphics using Unicode symbols
- Collision detection: players die on walls or each other
- Live scoreboard streamed to all players concurrently

## Usage
```bash
# Start server
go run cmd/server/main.go

# Connect players
go run cmd/client/main.go
```

## Controls
- **WASD** or **HJKL** to move

## Dependencies
- `github.com/gorilla/websocket`
