# Rabbits

A 2D multiplayer space game built with Go and Ebitengine. Control rabbit characters in a dynamic space environment with shooting mechanics, multiplayer support, and WebAssembly deployment for web browsers.

## Features

- **Multiple Game Modes**: Direct play, client-server multiplayer, and special "stars" mode
- **Real-time Multiplayer**: WebSocket-based networking for seamless multiplayer experience
- **Cross-platform**: Native desktop application and web browser support via WebAssembly
- **Dynamic Gameplay**: Movement, rotation, shooting mechanics with heat/load system
- **Scalable Architecture**: Scene-based system supporting different game types

## Quick Start

### Prerequisites

- **Go 1.21 or later**: [Download and install Go](https://golang.org/dl/)
- **Git**: For cloning the repository

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/demonodojo/rabbits.git
   cd rabbits
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the game:
   ```bash
   go run .
   ```

## Game Modes

### Menu Principal
Al ejecutar el juego sin argumentos, se mostrará un menú principal que permite elegir entre los diferentes modos de juego:

```bash
go run .
```

El menú incluye las siguientes opciones:
- **Juego Directo (Single Player)**: Modo de un solo jugador
- **Modo Stars**: Modo alternativo con mecánicas diferentes
- **Iniciar Servidor**: Inicia un servidor multiplayer en el puerto 8080
- **Conectar como Cliente**: Se conecta a un servidor en localhost:8080
- **Salir**: Cierra el juego

### Controles del Menú
- **Flechas ↑↓**: Navegar entre opciones
- **Enter o Espacio**: Seleccionar opción

### Modos de Juego Específicos

#### Single Player (Direct Mode)
```bash
go run . -direct
```

#### Multiplayer

##### Start a Server
```bash
go run . -server
```
The server will start on port 8080 and accept WebSocket connections.

##### Connect as Client
```bash
go run . -client
```
Connects to the server running on `localhost:8080`.

#### Stars Mode
```bash
go run . -stars
```
Alternative game mode with different mechanics.

## Controls

- **Arrow Keys**:
  - Left/Right: Rotate rabbit
  - Up/Down: Accelerate/Decelerate
- **F Key**: Fire (requires heat/load management)
- **Space**: Shoot (in some modes)

## Web Deployment

### Build for Web
```bash
make compile
```
This creates:
- `rabbits.wasm`: WebAssembly binary
- `wasm_exec.js`: WebAssembly runtime

### Serve Locally
Use any static file server to serve the web files:
```bash
# Using Python
python -m http.server 8000

# Using Node.js serve
npx serve .

# Using Go
go run -m http.server
```

### Docker Deployment
```bash
docker-compose up
```
Serves the game via nginx on port 80.

## Development

### Project Structure
```
├── game/              # Core game logic
│   ├── network/       # WebSocket client/server
│   ├── elements/      # UI components and forms
│   ├── scenes/        # Different game scenes
│   └── *.go          # Game entities (rabbit, player, bullet, etc.)
├── assets/           # Game assets (images, fonts)
├── web/              # Web deployment files
├── main.go           # Native application entry point
├── main_wasm.go      # WebAssembly entry point
└── Makefile          # Build scripts
```

### Build Commands

#### Manual WebAssembly Build
```bash
env NODE_ENV=development GOOS=js GOARCH=wasm go build -o rabbits.wasm
cp $(go env GOROOT)/misc/wasm/wasm_exec.js .
```

#### Development Tools
```bash
# Run tests
go test ./...

# Check for issues
go vet ./...

# Format code
go fmt ./...
```

## Network Configuration

### Server Configuration
- Default port: `:8080`
- WebSocket endpoint: `/ws`
- Accepts connections from any origin

### Client Configuration
- **Native**: Connects to `ws://localhost:8080/ws`
- **WebAssembly**: Connects to `ws://192.168.1.45:8080/ws` (update as needed)

## Architecture

The game uses a scene-based architecture built on Ebitengine:

- **Game**: Main game loop implementing `ebiten.Game`
- **Scene Interface**: Pluggable scenes for different game modes
- **Entity System**: Rabbit, Player, Bullet, Meteor entities with collision detection
- **Network Layer**: WebSocket-based with message queuing for smooth multiplayer
- **Asset Management**: Embedded assets for easy distribution

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and test thoroughly
4. Submit a pull request

## License

This project is licensed under the terms specified in the LICENSE file.

## Acknowledgments

- Built with [Ebitengine](https://ebitengine.org/) - A dead simple 2D game library for Go
- Inspired by classic space shooter games
- Originally based on the meteors example from Three Dots Labs
