package network

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Peer representa a un par conectado al servidor.
type Peer struct {
	Conn        *websocket.Conn
	IncomingMsg *MessageQueue
	OutgoingMsg *MessageQueue
	done        chan struct{}          // Canal para señalizar el cierre
	events      chan<- *websocket.Conn // Canal para publicar eventos de mensajes
}

func NewPeer(conn *websocket.Conn, events chan<- *websocket.Conn) *Peer {
	peer := &Peer{
		Conn:        conn,
		IncomingMsg: NewMessageQueue(),
		OutgoingMsg: NewMessageQueue(),
		done:        make(chan struct{}),
		events:      events,
	}
	go peer.readPump()
	go peer.writePump()
	return peer
}

func (p *Peer) readPump() {
	defer func() {
		p.Conn.Close()
	}()
	for {
		select {
		case <-p.done:
			return // Termina la goroutine si se recibe señal de cierre
		default:
			_, message, err := p.Conn.ReadMessage()
			if err != nil {
				// Manejar error o desconexión
				return
			}
			log.Printf("Mensaje recibido %s\n", string(message))
			p.IncomingMsg.Enqueue(string(message))
			p.events <- p.Conn
		}
	}
}

func (p *Peer) writePump() {
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-p.done:
			log.Printf("cerrando la cola...")
			return // Termina la goroutine si se recibe señal de cierre
		case <-ticker.C:
			message, ok := p.OutgoingMsg.Dequeue()
			if !ok {
				continue
			} else {
				log.Printf("Mensaje enviado %s\n", string(message))
				if err := p.Conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					// Manejar error
					return
				}
			}
		}
	}
}

func (p *Peer) Close() {
	close(p.done) // Cierra el canal para señalizar a las goroutines que terminen
}

func (p *Peer) Write(message string) {
	p.OutgoingMsg.Enqueue(message)
}

func (p *Peer) Read() (string, bool) {
	return p.IncomingMsg.Dequeue()
}
