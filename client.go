package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	connection *websocket.Conn // the websocket connection

	manager *Manager // reference to its manager/supervisor
	id      int      // identification from manager

	egress chan []byte // avoid concurrent writes by blocking a non-buffer channel
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

// client: read go routine => readMessages from outside/client
func (c *Client) readMessages() {
	defer func() {
		// graceful close the Conn once this process is done
		c.manager.removeClient(c)
	}()

	// loop forever => always runs as a go routine
	for {
		// ReadMessage is used to read the next message in queue of the Conn
		messageType, payload, err := c.connection.ReadMessage()

		// if error occurs, kill the Conn
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		log.Printf("client[%d] ReadMessage => "+"MessageType: %d\t"+"Payload: %s\n",
			c.id, messageType, string(payload))

		// Hack to test that WriteMessages works as intended
		// FIX: Will be replaced soon
		for wsClient := range c.manager.clients {
			wsClient.egress <- payload
		}
	}
}

// client: write go routine => writeMessages listens message to outside/client
func (c *Client) writeMessages() {
	defer func() {
		// graceful close the Conn once this process is done
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:

			// has been closed
			if !ok {
				// inform client of the close
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				// kill the go routine
				return
			}

			// write a regular text message to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("client[%d] WriteMessage: fail =>\n", err)
			} else {
				log.Printf("client[%d] WriteMessage: success\n", c.id)
			}
		}
	}
}
