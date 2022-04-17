package bootstrap

import (
	"context"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/pubsub"
	"strings"

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

	return srv.Run(ctx)
}
