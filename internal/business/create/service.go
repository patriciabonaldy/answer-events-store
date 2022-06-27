package create

import (
	"context"
	"encoding/json"

	"github.com/patriciabonaldy/bequest_challenge/internal/platform/pubsub"

	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"

	"github.com/patriciabonaldy/bequest_challenge/internal"
)

type Service interface {
	CreateAnswer(ctx context.Context, id string, data map[string]string) (*internal.Answer, error)
}

type service struct {
	producer   pubsub.Producer
	repository internal.Storage
	log        logger.Logger
}

// NewService returns the default Service interface implementation.
func NewService(producer pubsub.Producer, repository internal.Storage, log logger.Logger) Service {
	return &service{
		producer:   producer,
		repository: repository,
		log:        log,
	}
}

// CreateAnswer implements Service interface.
func (s service) CreateAnswer(ctx context.Context, id string, data map[string]string) (*internal.Answer, error) {
	body, err := json.Marshal(data)
	if err != nil {
		s.log.Errorf("error Marshal data %s", err.Error())
		return nil, err
	}

	event := internal.NewEvent(id, internal.Create, body)
	answer := internal.NewAnswer(event)
	m := generateMessage(answer)
	err = s.producer.Produce(ctx, m)
	if err != nil {
		return nil, err
	}

	return &answer, nil
}

func generateMessage(answer internal.Answer) pubsub.Message {
	message := pubsub.NewSystemMessage()
	message.ID = answer.ID
	message.CreateAt = answer.CreateAt
	message.UpdateAt = answer.UpdateAt

	for _, ev := range answer.Events {
		message.Events = append(message.Events, pubsub.Event{
			EventID:   ev.EventID,
			Type:      string(ev.Type),
			RawData:   ev.RawData,
			Timestamp: ev.Timestamp,
			Version:   ev.Version,
		})
	}

	return message
}
