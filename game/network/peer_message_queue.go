package network

import (
	"github.com/gorilla/websocket"
	"sync"
)

type PeerMessage struct {
	Peer    *websocket.Conn
	Message string
}

type PeerMessageQueue struct {
	mutex sync.Mutex
	items []PeerMessage
}

func NewPeerMessageQueue() *PeerMessageQueue {
	return &PeerMessageQueue{}
}

func (q *PeerMessageQueue) Enqueue(item PeerMessage) {
	q.mutex.Lock()
	q.items = append(q.items, item)
	q.mutex.Unlock()
}

func (q *PeerMessageQueue) Dequeue() (PeerMessage, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.items) == 0 {
		return PeerMessage{}, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *PeerMessageQueue) ReadAll() []PeerMessage {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	items := q.items
	q.items = nil // Vacía la cola
	return items
}
