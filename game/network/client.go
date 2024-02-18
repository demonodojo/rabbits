package network

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// Client representa a un cliente conectado a un servidor WebSocket.
type Client struct {
	Conn        *websocket.Conn
	IncomingMsg *MessageQueue
	OutgoingMsg *MessageQueue
	done        chan struct{}
}

func NewClient(url string) (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	client := &Client{
		Conn:        conn,
		IncomingMsg: NewMessageQueue(),
		OutgoingMsg: NewMessageQueue(),
		done:        make(chan struct{}),
	}
	go client.readPump()
	go client.writePump()
	return client, nil
}

func (c *Client) readPump() {
	defer close(c.done)
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		c.IncomingMsg.Enqueue(string(message))
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			message, ok := c.OutgoingMsg.Dequeue()
			if !ok {
				continue
			}
			log.Printf("Escribiendo: %s\n", message)
			if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func (c *Client) Close() {
	close(c.done)
	c.Conn.Close()
}

func (c *Client) Write(message string) {
	c.OutgoingMsg.Enqueue(message)
}

func (c *Client) Read() (string, bool) {
	return c.IncomingMsg.Dequeue()
}
