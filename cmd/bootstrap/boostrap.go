package bootstrap

import (
	"context"
	"strings"

	"github.com/patriciabonaldy/bequest_challenge/cmd/bootstrap/cdc"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/pubsub"

	"github.com/patriciabonaldy/big_queue/pkg/kafka"

	"github.com/patriciabonaldy/bequest_challenge/cmd/bootstrap/handler"
	"github.com/patriciabonaldy/bequest_challenge/internal/business"
	"github.com/patriciabonaldy/bequest_challenge/internal/config"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/storage/mongo"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	log := logger.New()
	ctx := context.Background()
	storage, err := mongo.NewDBStorage(ctx, cfg.Database, log)
	if err != nil {
		panic(err)
	}

	publisher := kafka.NewPublisher(strings.Split(cfg.Kafka.Broker, ","), cfg.Kafka.Topic)
	producer := pubsub.NewProducer(publisher)
	svc := business.NewService(producer, storage, log)
	handler := handler.New(svc, log)
	ctx, srv := New(ctx, cfg, handler)

	consumer := kafka.NewConsumer(strings.Split(cfg.Kafka.Broker, ","), cfg.Kafka.Topic)
	subscriber := pubsub.NewSubscriber(consumer, log)
	sync := cdc.NewSync(subscriber, storage, log)

	go func() {
		err = sync.Start(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	return srv.Run(ctx)
}
