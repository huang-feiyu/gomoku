package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	addr = "localhost"
	port = ":8080"

	// websocketUpgrader: incoming HTTP requests -> persitent WebSocket connection
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Manager holds references to all Clients Registered, and Broadcasting etc
type Manager struct {
}

// NewManager initalizes all the values inside the manager
func NewManager() *Manager {
	return &Manager{}
}

// serveWS: a HTTP Handler has the Manager that allows connections
func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")

	// upgrade the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// do nothing for now
	conn.Close()
	log.Println("Close connection")
}

// setupAPI will start all Routes and their Handlers
func setupAPI() {
	// create a Manager instance to handle WebSocket Connections
	manager := NewManager()

	// serve the ./client directory at Route /
	http.Handle("/", http.FileServer(http.Dir("./client")))

	http.HandleFunc("/ws", manager.serveWS)
}

func main() {
	setupAPI()

	// Server on port :8080
	log.Printf("Listen at: " + addr + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
