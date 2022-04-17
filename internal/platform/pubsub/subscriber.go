package pubsub

import (
	"context"
	"log"

	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/patriciabonaldy/big_queue/pkg"
)

type subscriber struct {
	consumer pkg.Consumer
	log      logger.Logger
}

type Subscriber interface {
	Subscriber(ctx context.Context, callback func(ctx context.Context, message interface{}) error) error
}

func NewSubscriber(consumer pkg.Consumer, log logger.Logger) Subscriber {
	p := subscriber{
		consumer: consumer,
		log:      log,
	}

	return &p
}

func (s subscriber) Subscriber(ctx context.Context, callback func(ctx context.Context, message interface{}) error) error {
	chMsg := make(chan pkg.Message)
	chErr := make(chan error)
	go func() {
		s.consumer.Read(ctx, chMsg, chErr)
	}()

	// read/process message
	for {
		select {
		case m := <-chMsg:
			log.Println(m)
			callback(ctx, m)
		case err := <-chErr:
			log.Println(err)
		}
	}
}
