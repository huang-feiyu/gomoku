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
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

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
