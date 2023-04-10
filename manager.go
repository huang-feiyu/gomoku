package main

import (
	"log"
	"net/http"
	"sync"
)

// Manager holds references to all Clients Registered, and Broadcasting etc
type Manager struct {
	id int // identification for client

	clients ClientList

	// Lock when updating clients
	sync.RWMutex
}

// NewManager initalizes all the values inside the manager
func NewManager() *Manager {
	return &Manager{
		id:      1,
		clients: make(ClientList),
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

	log.Printf("New connection: client[%d]\n", client.id)

}

// addClient will add clients to our clientList
func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	client.id = m.id
	m.id++

	m.clients[client] = true
}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		_ = client.connection.Close()
		delete(m.clients, client)

		log.Println("Close connection: client[", client.id, "]")
	}
}
