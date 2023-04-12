package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var (
	pongWait     = 10 * time.Second    // pongWait is how long we will await a pong response from client
	pingInterval = (pongWait * 9) / 10 // pingInterval has to be less than pongWait, otherwise server will close before next ping

	ROLE_UNMATCH = 0
	ROLE_PALYER1 = 1
	ROLE_PALYER2 = 2
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	connection *websocket.Conn // the websocket connection

	manager *Manager // reference to its manager/supervisor
	id      int      // identification from manager

	egress chan Event // avoid concurrent writes by blocking a non-buffer channel

	room *Room // group field
	role int
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		room:       &Room{},
		role:       ROLE_UNMATCH,
	}
}

// GetPartner returns the opponent of this client
func (c *Client) GetPartner() *Client {
	for client := range c.manager.clients {
		if client.room.id == c.room.id && (client.role == ROLE_PALYER1 && c.role == ROLE_PALYER2 ||
			client.role == ROLE_PALYER2 && c.role == ROLE_PALYER1) {
			return client
		}
	}
	return nil
}

// client: read go routine => readMessages from outside/client
func (c *Client) readMessages() {
	defer func() {
		// graceful close the Conn once this process is done
		c.manager.removeClient(c)
	}()
	c.connection.SetReadLimit(512)

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	c.connection.SetPongHandler(c.pongHandler)

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

// pongHandler is used to handle PongMessages from Client
func (c *Client) pongHandler(pongMsg string) error {
	log.Printf("client[%d] Pong: return to server\n", c.id)
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

// client: write go routine => writeMessages listens message to outside/client
func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
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
					log.Printf("client[%d] WriteMessage: already closed connection => %v", c.id, err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			// write a regular text message to the connection
			if err = c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("client[%d]", c.id)
				log.Println("WriteMessage: fail => ", err)
			} else {
				log.Printf("client[%d] WriteMessage: sent(%s) to client\n", c.id, message.Type)
			}

		case <-ticker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("client[%d] Ping: fail => %v", c.id, err)
				return
			}
			log.Printf("client[%d] Ping: sent to client\n", c.id)
		}

	}

}
