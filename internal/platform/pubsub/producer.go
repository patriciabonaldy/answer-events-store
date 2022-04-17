package pubsub

import (
	"context"
	"fmt"

	"github.com/patriciabonaldy/big_queue/pkg"
)

type producer struct {
	publisher pkg.Publisher
}

type Producer interface {
	Produce(ctx context.Context, event interface{}) error
}

func NewProducer(publisher pkg.Publisher) Producer {
	p := producer{
		publisher: publisher,
	}

	return &p
}

func (p producer) Produce(ctx context.Context, event interface{}) error {
	if err := p.publisher.Publish(ctx, event); err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}
	return nil
}
