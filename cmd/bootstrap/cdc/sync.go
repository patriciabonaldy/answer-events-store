package cdc

import (
	"context"
	"encoding/json"

	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/pubsub"
)

type service struct {
	subscriber pubsub.Subscriber
	repository internal.Storage
	log        logger.Logger
}

func NewSync(subscriber pubsub.Subscriber, repository internal.Storage, log logger.Logger) *service {
	return &service{
		subscriber: subscriber,
		repository: repository,
		log:        log,
	}
}

func (s *service) Start(ctx context.Context) error {
	return s.subscriber.Subscriber(ctx, s.callBack)
}

func (s *service) callBack(ctx context.Context, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		s.log.Error("invalid message type")
		return err
	}

	var msg pubsub.Message
	err = json.Unmarshal(data, &msg)
	if err != nil {
		s.log.Error("invalid message type")
		return err
	}

	answer := toAnswer(msg)
	ev := answer.Events[len(answer.Events)-1]

	switch {
	case ev.Type == internal.Create && len(answer.Events) == 1:
		err = s.repository.Save(ctx, answer)
		if err != nil {
			s.log.Errorf("error CreateAnswer %s", err.Error())
			return err
		}

	default:
		err = s.repository.Update(ctx, answer)
		if err != nil {
			s.log.Errorf("error UpdateAnswer ID:%s:%s", err.Error())
			return err
		}
	}

	return nil
}

func toAnswer(message pubsub.Message) internal.Answer {
	answer := internal.Answer{
		ID:       message.ID,
		CreateAt: message.CreateAt,
		UpdateAt: message.UpdateAt,
	}

	for _, ev := range message.Events {
		answer.Events = append(answer.Events, internal.Event{
			EventID:   ev.EventID,
			Type:      internal.EventType(ev.Type),
			RawData:   ev.RawData,
			Timestamp: ev.Timestamp,
			Version:   ev.Version,
		})
	}

	return answer
}
