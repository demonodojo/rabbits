//go:build js && wasm

package network

import (
	"context"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"syscall/js"
)

// JSClient representa a un cliente conectado a un servidor WebSocket utilizando syscall/js.
type JSClient struct {
	Ws          *websocket.Conn
	IncomingMsg *MessageQueue
	OutgoingMsg *MessageQueue
	connected   bool
	done        chan struct{}
	newMessage  chan struct{}
}

func NewJSClient(url string) (*JSClient, error) {
	client := &JSClient{
		IncomingMsg: NewMessageQueue(),
		OutgoingMsg: NewMessageQueue(),
		connected:   true,
		done:        make(chan struct{}),
		newMessage:  make(chan struct{}, 1), // No bloqueante
	}
	c, _, err := websocket.Dial(context.Background(), url, nil)
	if err != nil {
		fmt.Println(err, "ERROR")
	}
	client.Ws = c

	go client.writePump()
	go client.readPump()
	return client, nil
}

// onOpen se dispara cuando la conexión WebSocket está abierta y lista para enviar mensajes.
func (c *JSClient) onOpen(this js.Value, args []js.Value) interface{} {
	c.connected = true
	println("WebSocket connection is open.")
	return nil
}

func (c *JSClient) readPump() {
	defer func() {
		// Close the WebSocket connection when the function returns.
		//c.Ws.Close(websocket.StatusGoingAway, "BYE")
	}()

	for {
		// Read a message from the WebSocket connection.
		messageType, payload, err := c.Ws.Read(context.Background())

		if err != nil {
			// Log and panic if there is an error reading the message.
			log.Panicf(err.Error())
		}
		c.IncomingMsg.Enqueue(string(payload))

		// Log the message type and payload for debugging.
		log.Println("MessageType: ", messageType)
		log.Println("Payload: ", string(payload))
	}
}

func (c *JSClient) writePump() {
	for {
		select {
		case <-c.done:
			return
		case <-c.newMessage:
			if c.connected {
				message, ok := c.OutgoingMsg.Dequeue()
				if !ok {
					continue
				}
				err := c.Ws.Write(context.Background(), websocket.MessageText, []byte(message))
				if err != nil {
					log.Println("Error writing to WebSocket:", err)
				}
			}
		}
	}
}

func (c *JSClient) Close() {
	close(c.done)
	c.Ws.Close(websocket.StatusGoingAway, "BYE")
}

func (c *JSClient) Write(message string) {
	println("writing message")
	c.OutgoingMsg.Enqueue(message)
	select {
	case c.newMessage <- struct{}{}:
	default:
	}
}

func (c *JSClient) Read() (string, bool) {
	return c.IncomingMsg.Dequeue()
}

func (c *JSClient) ReadAll() []string {
	return c.IncomingMsg.ReadAll()
}
