package main

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

// Manager holds references to all Clients Registered, and Broadcasting etc
type Manager struct {
	clients ClientList

	id int // identification for client

	room int // identification for room

	sync.RWMutex

	handlers map[string]EventHandler // handlers: map[event_type] -> handler
}

// NewManager initializes all the values inside the manager
func NewManager() *Manager {
	m := &Manager{
		id:       1,
		room:     1,
		clients:  make(ClientList),
		handlers: map[string]EventHandler{},
	}
	m.setupEventHandlers()
	return m
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventChangeName] = ChangeNameHandler
}

// routeEvent makes sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

// serveWS: HTTP Handler has the Manager that allows connections
func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	// upgrade the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// create new client & add to manager
	client := NewClient(conn, m)
	m.addClient(client)

	// start two go routines
	go client.readMessages()
	go client.writeMessages()

	log.Printf("client[%d] New connection: starts to read/write\n", client.id)

	// match if possible
	m.matchClient(client)
}

// addClient will add clients to our clientList
func (m *Manager) addClient(c *Client) {
	m.Lock()
	defer m.Unlock()

	c.id = m.id
	m.id++

	m.clients[c] = true
}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		_ = client.connection.Close()
		delete(m.clients, client)

		log.Printf("client[%d] Close connection\n", client.id)
	}
}

// matchClient will match if possible
func (m *Manager) matchClient(c *Client) {
	m.Lock()
	defer m.Unlock()

	// match if possible
	for client := range m.clients {
		if client.role == ROLE_UNMATCH && c != client {
			client.room = m.room
			c.room = m.room
			m.room++

			_ = SendConnectRole(client, ROLE_PALYER1)
			client.role = ROLE_PALYER1

			_ = SendConnectRole(c, ROLE_PALYER2)
			c.role = ROLE_PALYER2

			log.Printf("client[%d] Match client[%d] at room[%d]\n", client.id, c.id, m.room-1)
			break
		}
	}
}
