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

var clientManager = NewClientManager()

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar WebSocket:", err)
		return
	}
	log.Println("Nueva Conexión")
	clientManager.RegisterClient(ws)

}

func (s *Server) Start() {

	go clientManager.Run()

	go func() {
		http.HandleFunc("/ws", s.handleConnections)
		log.Printf("Iniciando servidor WebSocket en %s\n", s.Port)
		if err := http.ListenAndServe(s.Port, nil); err != nil {
			log.Fatal("Error iniciando servidor WebSocket:", err)
		}
	}()
}

func (s *Server) ReadAll() []PeerMessage {
	return clientManager.allMessages.ReadAll()
}

func (s *Server) Broadcast(message string) {
	clientManager.Broadcast(message)
}
