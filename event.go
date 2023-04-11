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
	EventChangeName  = "name_message"
	EventMove        = "move_message"
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
		if client.room.id == c.room.id {
			client.egress <- outgoingEvent
		}
	}
	return nil
}

type ConnectRoleEvent struct {
	Role int `json:"role"`
}

// SendConnectRole is called when a client emerges
// AND this client match with the waiting client => send role message to client
func SendConnectRole(client *Client, role int) error {
	playerEvent := ConnectRoleEvent{role}
	data, _ := json.Marshal(playerEvent)
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventConnectRole
	client.egress <- outgoingEvent

	return nil
}

type ChangeNameEvent struct {
	Name string `json:"name"`
	Role int    `json:"role"`
}

// ChangeNameHandler not only changes the current client name, => receive name from client
// but also updates the pair's display => send name message to the pair
func ChangeNameHandler(event Event, c *Client) error {
	// receive name from client
	var changeNameEvent ChangeNameEvent
	if err := json.Unmarshal(event.Payload, &changeNameEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// send name to opponent as well as THIS client

	data, _ := json.Marshal(changeNameEvent)
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventChangeName
	opponent := c.GetPartner()
	if opponent == nil {
		return fmt.Errorf("change name is not allowed if no partner")
	}
	opponent.egress <- outgoingEvent
	c.egress <- outgoingEvent
	return nil
}

type MoveEvent struct {
	Role int `json:"role"`
	Row  int `json:"row"`
	Col  int `json:"col"`
}

// MoveHandler not only receive the move, => receive move from client
// but also updates the pair's display => send move message to the pair
// TODO: result
func MoveHandler(event Event, c *Client) error {
	// receive move from client
	var moveEvent MoveEvent
	if err := json.Unmarshal(event.Payload, &moveEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	res := c.room.Move(moveEvent.Role, moveEvent.Row, moveEvent.Col)

	// TODO: send display to both of the pair

	// send result to the pair
	if res != 0 {
	}
	return nil
}
