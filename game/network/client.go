package network

import (
	"github.com/gorilla/websocket"
	"log"
)

func connectToServer() {
	// Conectar al servidor WebSocket
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Enviar un mensaje al servidor
	err = c.WriteMessage(websocket.TextMessage, []byte("Hola desde el cliente Ebiten!"))
	if err != nil {
		log.Println("write:", err)
		return
	}

	// Leer respuesta (Si aplicable)
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Printf("Recibido: %s\n", message)
}

func main() {
	// Aquí inicias tu juego Ebiten y te conectas al servidor WebSocket
	connectToServer()
}
