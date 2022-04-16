package business

import (
	"context"
	"encoding/json"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"time"

	"github.com/patriciabonaldy/bequest_challenge/internal"
)

var (
	mapStatusValid = map[internal.EventType]map[internal.EventType]interface{}{
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
	CreateAnswer(ctx context.Context, data map[string]string) (*internal.Answer, error)
	GetAnswerByID(ctx context.Context, eventID string) (*internal.Answer, error)
	UpdateAnswer(ctx context.Context, eventID, eventType string, data map[string]string) error
}

type service struct {
	repository internal.Storage
	log        logger.Logger
}

// NewService returns the default Service interface implementation.
func NewService(repository internal.Storage, log logger.Logger) Service {
	return &service{
		repository: repository,
		log:        log,
	}
}

// CreateAnswer implements Service interface.
func (s service) CreateAnswer(ctx context.Context, data map[string]string) (*internal.Answer, error) {
	body, err := json.Marshal(data)
	if err != nil {
		s.log.Errorf("error Marshal data %s", err.Error())
		return nil, err
	}

	event := internal.NewEvent("", internal.Create, body)
	answer := internal.NewAnswer(event)
	err = s.repository.Save(ctx, answer)
	if err != nil {
		s.log.Errorf("error CreateAnswer %s", err.Error())
		return nil, err
	}

	return &answer, nil
}

func (s service) GetAnswerByID(ctx context.Context, eventID string) (*internal.Answer, error) {
	answer, err := s.repository.GetByID(ctx, eventID)
	if err != nil {
		s.log.Errorf("error GetAnswerByID ID:%s:%s", eventID, err.Error())
		return nil, err
	}

	return &answer, nil
}

// UpdateAnswer implements Service interface.
func (s service) UpdateAnswer(ctx context.Context, eventID, eventType string, data map[string]string) error {
	answer, err := s.repository.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if err = s.checkEvent(answer, eventType); err != nil {
		s.log.Errorf("error checkEvent ID:%s:%s", eventID, err.Error())
		return err
	}

	body := answer.Events[len(answer.Events)-1].RawData
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			s.log.Errorf("error Marshal data %s", err.Error())
			return err
		}
	}

	event := internal.NewEvent("", internal.EventType(eventType), body)
	answer.AddEvent(event)
	answer.UpdateAt = time.Now()
	err = s.repository.Update(ctx, answer)
	if err != nil {
		s.log.Errorf("error UpdateAnswer ID:%s:%s", err.Error())
		return err
	}

	return nil
}

func (s service) checkEvent(answer internal.Answer, eventType string) error {
	event := answer.Events[len(answer.Events)-1].Type
	switch event {
	case internal.Create:
		_, ok := mapStatusValid[internal.Create][internal.EventType(eventType)]
		if !ok {
			return internal.ErrInvalidEventStatus
		}
	case internal.Delete:
		_, ok := mapStatusValid[internal.Delete][internal.EventType(eventType)]
		if !ok {
			return internal.ErrInvalidEventStatus
		}
	case internal.Update:
		_, ok := mapStatusValid[internal.Update][internal.EventType(eventType)]
		if !ok {
			return internal.ErrInvalidEventStatus
		}

	default:
		return internal.ErrInvalidEventStatus
	}

	return nil
}
