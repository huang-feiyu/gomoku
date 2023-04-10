package main

import (
	"log"
	"net/http"
)

var (
	port = ":8080"
)

func main() {
	setupAPI()

	// Server on port :8080
	log.Printf("Listen at: localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// setupAPI will start all Routes and their Handlers
func setupAPI() {
	// Serve the ./client directory at Route /
	http.Handle("/", http.FileServer(http.Dir("./client")))
}
