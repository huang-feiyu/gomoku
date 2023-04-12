package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	addr = "localhost"
	port = ":8080"

	// websocketUpgrader: incoming HTTP requests -> persitent WebSocket connection
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// checkOrigin will check origin and return true if its allowed
func checkOrigin(r *http.Request) bool {
	return true
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
