package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type Storage interface {
	GetByID(ctx context.Context, ID string) (*Answer, error)
	Save(ctx context.Context, answer Answer) error
	Update(ctx context.Context, answer Answer) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=platform/storage/storagemocks --name=Storage

// Answer is a structure of answers to be stored
type Answer struct {
	AnswerID primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Events   []Event             `bson:"events,omitempty"`
	CreateAt primitive.Timestamp `json:"createdAt" bson:"createdAt,omitempty"`
	UpdateAt primitive.Timestamp `json:"updatedAt" bson:"updatedAt,omitempty"`
}

// Event is a structure of events to be stored
type Event struct {
	EventID   primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Type      string              `bson:"event_type"`
	RawData   bson.Raw            `bson:"data,omitempty"`
	Timestamp primitive.Timestamp `timestamp:"createdAt" bson:"timestamp"`
	Version   int                 `bson:"version"`
}

func NewEvent(eventType string, data []byte, version int) (Event, error) {
	createdAt := primitive.Timestamp{
		T: uint32(time.Now().Unix()),
	}

	return Event{
		Type:      eventType,
		RawData:   bson.Raw{Kind: data[0], Data: data},
		Timestamp: createdAt,
		Version:   version,
	}, nil
}

func (c *Event) ID() string {
	return c.EventID.String()
}

func (c *Event) Name() string {
	return c.Type
}

func (c *Event) At() time.Time {
	return time.Unix(int64(c.Timestamp.T)/1000, int64(c.Timestamp.T)%1000*1000000)
}

func (c *Event) Data() []byte {
	return c.RawData.Data
}

func (c *Event) Unmarshall(out interface{}) error {
	return c.RawData.Unmarshal(out)
}
