package internal

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type Storage interface {
	GetByID(ctx context.Context, ID string) (Answer, error)
	Save(ctx context.Context, answer Answer) error
	Update(ctx context.Context, answer Answer) error
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=platform/storage/storagemocks --name=Storage

type EventType string

var (
	Create EventType = "create"
	Update EventType = "update"
	Delete EventType = "delete"
)

var (
	ErrInvalidEvent = errors.New("event can not be empty")
	ErrInvalidData  = errors.New("data can not be empty")

	ErrIDIsEmpty          = errors.New("invalid ID")
	ErrAnswerNotFound     = errors.New("id not found")
	ErrInvalidEventStatus = errors.New("invalid event")
)

// Answer is a structure of answers to be stored
type Answer struct {
	ID       string
	Events   []Event
	CreateAt time.Time
	UpdateAt time.Time
}

// Event is a structure of events to be stored
type Event struct {
	EventID   string
	Type      EventType
	RawData   []byte
	Timestamp time.Time
	Version   int
}

func NewAnswer(event Event) Answer {
	id, _ := uuid.NewUUID()

	return Answer{
		ID:       id.String(),
		Events:   []Event{event},
		CreateAt: time.Now(),
	}
}

func NewEvent(eventID string, eventType EventType, data []byte) Event {
	if eventID == "" {
		id, _ := uuid.NewUUID()
		eventID = id.String()
	}

	return Event{
		EventID:   eventID,
		Type:      eventType,
		RawData:   data,
		Timestamp: time.Now(),
	}
}

func (a *Answer) AddEvent(event Event) {
	event.Version = len(a.Events)
	a.Events = append(a.Events, event)
}
