package internal

import (
	"context"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Answer is a structure of answers to be stored
type Answer struct {
	ID       string    `bson:"_id"`
	Events   []Event   `bson:"events"`
	CreateAt time.Time `bson:"create_at"`
	UpdateAt time.Time `bson:"update_at"`
}

// Event is a structure of events to be stored
type Event struct {
	Type      string    `bson:"event_type"`
	RawData   bson.Raw  `bson:"data,omitempty"`
	Timestamp time.Time `bson:"timestamp"`
	Version   int       `bson:"version"`
}

type Storage interface {
	GetByID(ctx context.Context, ID string) (*Answer, error)
	Save(ctx context.Context, answer Answer) error
	Update(ctx context.Context, answer Answer) error
}
