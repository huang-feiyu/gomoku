package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Event is the Messages sent over the WebSocket
// Used to differ between different actions
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// EventHandler interface
type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
	EventConnectRole = "role_message" // when connect
)

// SendMessageEvent is the payload sent in the send_message event
type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

// NewMessageEvent is returned when responding to send_message
type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

// SendMessageHandler will send out a message to all other participants in the chat
func SendMessageHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// Prepare an Outgoing Message to others
	var broadMessage NewMessageEvent

	broadMessage.Sent = time.Now()
	broadMessage.Message = chatevent.Message
	broadMessage.From = chatevent.From

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Place payload into an Event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage
	// Broadcast to all other Clients in the same chat group
	for client := range c.manager.clients {
		if client.room == c.room {
			client.egress <- outgoingEvent
		}
	}
	return nil
}

type ConnectRoleEvent struct {
	Role int `json:"role"`
}

func SendConnectRole(client *Client, role int) error {
	playerEvent := ConnectRoleEvent{role}
	data, err := json.Marshal(playerEvent)
	if err != nil {
		return err
	}
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventConnectRole
	client.egress <- outgoingEvent

	return nil
}
