# Socky

 Multiplayer snake game inside the terminal, written in Go using WebSockets. Server handles game logic, Clients send keypresses.

**Controls for client:** `wasd` or `hjkl` (vim keys)

**Run:**
```bash
go run main.go    # server
go run client.go  # client
```


## Dependencies
- `github.com/gorilla/websocket`
