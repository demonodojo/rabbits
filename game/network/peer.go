package network

import (
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
			p.IncomingMsg.Enqueue(string(message))
			p.events <- p.Conn
		}
	}
}

func (p *Peer) writePump() {
	for {
		select {
		case <-p.done:
			return // Termina la goroutine si se recibe señal de cierre
		default:
			message, ok := p.OutgoingMsg.Dequeue()
			if !ok {
				// Esperar o manejar cuando no hay mensajes
				return
			}
			if err := p.Conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				// Manejar error
				return
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
