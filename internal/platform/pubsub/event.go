package pubsub

import (
	"time"
)

type Message struct {
	ID       string    `json:"message_id"`
	Events   []Event   `json:"events"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

// Event request model.
type Event struct {
	EventID   string    `json:"event_id"`
	Type      string    `json:"type"`
	RawData   []byte    `json:"data"`
	Timestamp time.Time `json:"at"`
	Version   int       `json:"version"`
}

// NewSystemMessage when we want to publish global message to all users in the room
func NewSystemMessage() Message {
	return Message{
		ID:       "",
		Events:   []Event{},
		CreateAt: time.Time{},
		UpdateAt: time.Time{},
	}
}
