package business

import (
	"context"

	"github.com/pkg/errors"

	"github.com/patriciabonaldy/bequest_challenge/internal"
)

var (
	ErrInvalidEventStatus = errors.New("invalid event")
	mapStatusValid        = map[internal.EventType]map[internal.EventType]interface{}{
		internal.Create: {
			internal.Delete: nil,
			internal.Update: nil,
		},
		internal.Delete: {
			internal.Create: nil,
		},
		internal.Update: {
			internal.Delete: nil,
		},
	}
)

type Service interface {
	CreateAnswer(ctx context.Context, data []byte) error
	GetAnswers(ctx context.Context, eventID string) (internal.Answer, error)
	UpdateAnswer(ctx context.Context, eventID, eventType string, data []byte) error
}

type service struct {
	repository internal.Storage
}

// NewService returns the default Service interface implementation.
func NewService(repository internal.Storage) Service {
	return &service{
		repository: repository,
	}
}

// CreateAnswer implements Service interface.
func (s service) CreateAnswer(ctx context.Context, data []byte) error {
	event := internal.NewEvent("", internal.Create, data)
	answer := internal.NewAnswer(event)

	_, err := s.repository.Save(ctx, answer)
	if err != nil {
		return err
	}

	return nil
}

// UpdateAnswer implements Service interface.
func (s service) UpdateAnswer(ctx context.Context, eventID, eventType string, data []byte) error {
	answer, err := s.repository.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if err = s.checkEvent(answer, eventType); err != nil {
		return err
	}

	event := internal.NewEvent("", internal.EventType(eventType), data)
	answer.AddEvent(event)
	answer, err = s.repository.Update(ctx, answer)
	if err != nil {
		return err
	}

	return nil
}

func (s service) GetAnswers(ctx context.Context, eventID string) (internal.Answer, error) {
	answer, err := s.repository.GetByID(ctx, eventID)
	if err != nil {
		return internal.Answer{}, err
	}

	return answer, nil
}

func (s service) checkEvent(answer internal.Answer, eventType string) error {
	event := answer.Events[len(answer.Events)-1].Type
	switch event {
	case internal.Create:
		_, ok := mapStatusValid[internal.Create][internal.EventType(eventType)]
		if !ok {
			return ErrInvalidEventStatus
		}
	case internal.Delete:
		_, ok := mapStatusValid[internal.Delete][internal.EventType(eventType)]
		if !ok {
			return ErrInvalidEventStatus
		}
	case internal.Update:
		_, ok := mapStatusValid[internal.Update][internal.EventType(eventType)]
		if !ok {
			return ErrInvalidEventStatus
		}

	default:
		return ErrInvalidEventStatus
	}

	return nil
}
