// client_manager.go

package network

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// ClientManager mantiene un registro de todas las conexiones de clientes WebSocket.
type ClientManager struct {
	peers         map[*websocket.Conn]*Peer // Un mapa para mantener un registro de las conexiones
	register      chan *websocket.Conn      // Un canal para registrar nuevas conexiones
	unregister    chan *websocket.Conn      // Un canal para desregistrar conexiones existentes
	messageEvents chan *websocket.Conn      // Cambiado para transportar solo el identificador del Peer

	allMessages *PeerMessageQueue
	mutex       sync.Mutex
}

// NewClientManager crea e inicializa una nueva instancia de ClientManager.
func NewClientManager() *ClientManager {
	return &ClientManager{
		peers:         make(map[*websocket.Conn]*Peer),
		register:      make(chan *websocket.Conn),
		unregister:    make(chan *websocket.Conn),
		messageEvents: make(chan *websocket.Conn),
		allMessages:   NewPeerMessageQueue(),
	}
}

// Run inicia el proceso de manejo de las conexiones de clientes.
func (manager *ClientManager) Run() {
	for {
		select {
		case conn := <-manager.register:
			manager.mutex.Lock()
			manager.peers[conn] = NewPeer(conn, manager.messageEvents)
			manager.mutex.Unlock()

		case conn := <-manager.unregister:
			manager.mutex.Lock()
			if peer, ok := manager.peers[conn]; ok {
				delete(manager.peers, conn)
				peer.Close() // Asegúrate de cerrar el Peer adecuadamente
			}
			manager.mutex.Unlock()
		case peerConn := <-manager.messageEvents:
			if peer, ok := manager.peers[peerConn]; ok {
				// Intenta leer un mensaje de la cola de salida del Peer
				log.Println("Aviso de mensaje entrante")
				if message, ok := peer.Read(); ok {
					// Procesa el mensaje, por ejemplo, encolándolo en allMessages
					log.Printf("Mensaje leido de la cola del peer %s", message)
					manager.allMessages.Enqueue(PeerMessage{Peer: peerConn, Message: message})
				} else {
					log.Println("Nada en la cola")
				}
			}
		}
	}
}

// RegisterClient añade una nueva conexión de cliente al ClientManager.
func (manager *ClientManager) RegisterClient(conn *websocket.Conn) {
	manager.register <- conn
}

// UnregisterClient elimina una conexión de cliente existente del ClientManager.
func (manager *ClientManager) UnregisterClient(conn *websocket.Conn) {
	manager.unregister <- conn
}

// GetClients devuelve una lista de todas las conexiones de clientes activas.
func (manager *ClientManager) GetClients() []*websocket.Conn {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	clients := make([]*websocket.Conn, 0, len(manager.peers))
	for conn := range manager.peers {
		clients = append(clients, conn)
	}
	return clients
}

func (manager *ClientManager) Close() {

	for conn := range manager.peers {
		manager.UnregisterClient(conn)
	}
}

func (manager *ClientManager) Write(conn *websocket.Conn, message string) {
	if peer, ok := manager.peers[conn]; ok {
		peer.Write(message)
	}
}

func (manager *ClientManager) Read(conn *websocket.Conn, message string) (string, bool) {
	if peer, ok := manager.peers[conn]; ok {
		return peer.Read()
	} else {
		return "", false
	}
}

func (manager *ClientManager) ReadAll() []PeerMessage {
	return manager.allMessages.ReadAll()
}

func (manager *ClientManager) Broadcast(message string) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for _, peer := range manager.peers {
		peer.Write(message)
	}
}
