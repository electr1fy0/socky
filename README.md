# Socky: Multiplayer Terminal Snake

Socky is a real-time, multiplayer snake game playable directly in your terminal. Challenge friends in private rooms, control your snake with simple keypresses, and experience fast-paced, dynamic gameplay.

## Screenshots
<div style="gap:10px;">
<img style="width: 60%; height="auto" alt="Screenshot 2025-09-18 at 11 34 50" src="https://github.com/user-attachments/assets/defde25f-ddb0-4ff0-a5e4-ba9ef4664f76" />
<img style="width: 60%; height="auto" alt="Screenshot 2025-09-18 at 11 41 03" src="https://github.com/user-attachments/assets/e89495ae-e4e8-4ef9-9828-6f398c3f6db6" />
</div>

## Features

- Real-time multiplayer gameplay
- Private game rooms
- Unicode-based terminal graphics
- Wall and snake collision handling
- Live, shared scoreboard
- Keyboard-driven controls
## Getting Started

### Prerequisites

-   [Go](https://golang.org/dl/) (Go 1.22+)

### Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/electr1fy0/socky.git
    cd socky
    ```

2.  Install dependencies:

    ```bash
    go mod tidy
    ```

### Usage

1.  **Start the Server:**

    ```bash
    go run cmd/server/main.go
    ```
    The server will start on `http://localhost:8080`.

2.  **Create a Game Room (from a client terminal):**

    ```bash
    go run cmd/client/main.go create
    ```
    This command will create a new room and print its `ROOM_ID`. Make a note of this ID.

3.  **Join a Game Room (from a client terminal):**

    ```bash
    go run cmd/client/main.go <ROOM_ID>
    ```
    Replace `<ROOM_ID>` with the ID obtained in the previous step. Multiple players can join the same room.

## Controls

-   **WASD** or **HJKL** to move
-   **Q** to quit the game

## Dependencies

-   `github.com/coder/websocket`
-   `github.com/google/uuid`
-   `golang.org/x/term`

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
