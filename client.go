package main

import (
	"encoding/json"
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

	egress chan Event // avoid concurrent writes by blocking a non-buffer channel
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
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
		_, payload, err := c.connection.ReadMessage()

		// if error occurs, kill the Conn
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		// Marshal incoming data into a Event struct
		var request Event
		if err = json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message: %v", err)
			break // Breaking the connection here might be harsh xD
		}
		// Route the Event
		if err = c.manager.routeEvent(request, c); err != nil {
			log.Println("Error handling Message: ", err)
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
					log.Printf("client[%d]", c.id)
					log.Println("WriteMessage: already closed connection => ", err)
				}
				return // kill the go routine
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return // kill the go routine
			}
			// write a regular text message to the connection
			if err = c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("client[%d]", c.id)
				log.Println("WriteMessage: fail => ", err)
			} else {
				log.Printf("client[%d] WriteMessage: success\n", c.id)
			}
		}
	}
}
