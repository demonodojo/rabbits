// game/network/queue.go

package network

import (
	"sync"
)

// MessageQueue representa una cola FIFO segura para concurrencia para almacenar mensajes.
type MessageQueue struct {
	mutex sync.Mutex
	items []string
}

// NewMessageQueue crea una nueva instancia de MessageQueue.
func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		items: make([]string, 0),
	}
}

// Enqueue agrega un elemento al final de la cola.
func (q *MessageQueue) Enqueue(item string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.items = append(q.items, item)
}

// Dequeue elimina y devuelve el primer elemento de la cola.
// Devuelve "", false si la cola está vacía.
func (q *MessageQueue) Dequeue() (string, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.items) == 0 {
		return "", false
	}

	item := q.items[0]
	q.items = q.items[1:] // Deslizar el slice para remover el elemento.
	return item, true
}

// Size devuelve el número actual de elementos en la cola.
func (q *MessageQueue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.items)
}
