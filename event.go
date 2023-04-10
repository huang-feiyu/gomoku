package main

import "encoding/json"

// Event is the Messages sent over the WebSocket
// Used to differ between different actions
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// EventHandler interface
type EventHandler func(event Event, c *Client) error

const (
	// EventSendMessage is the event type for new chat messages sent
	EventSendMessage = "send_message"
)

// SendMessageEvent is the payload sent in the send_message event
type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}
