// network/server.go

package network

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Server struct {
	Port string // Dirección en la que el servidor escuchará
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Aceptar cualquier origen por simplicidad
	},
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error al leer mensaje:", err)
			break
		}
		log.Printf("Mensaje recibido: %s\n", message)

		// Maneja aquí el mensaje recibido
	}
}

func (s *Server) Start() {
	http.HandleFunc("/ws", s.handleConnections)
	log.Printf("Iniciando servidor WebSocket en %s\n", s.Port)
	if err := http.ListenAndServe(s.Port, nil); err != nil {
		log.Fatal("Error iniciando servidor WebSocket:", err)
	}
}
