package bootstrap

import (
	"context"
	"strings"

	"github.com/patriciabonaldy/bequest_challenge/internal/business/create"
	"github.com/patriciabonaldy/bequest_challenge/internal/business/find"
	"github.com/patriciabonaldy/bequest_challenge/internal/business/update"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/bus/inmemory"

	"github.com/patriciabonaldy/bequest_challenge/cmd/bootstrap/cdc"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/pubsub"

	"github.com/patriciabonaldy/big_queue/pkg/kafka"

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

	publisher := kafka.NewPublisher(strings.Split(cfg.Kafka.Producer, ","), cfg.Kafka.Topic)
	producer := pubsub.NewProducer(publisher)
	commandBus := inmemory.NewCommandBus()

	// create handler
	createSvc := create.NewService(producer, storage, log)
	createAnswerCommandHandler := create.NewAnswerCommandHandler(createSvc)
	commandBus.Register(create.AnswerCommandType, createAnswerCommandHandler)

	// find handler
	findSvc := find.NewService(storage, log)
	findAnswerCommandHandler := find.NewAnswerCommandHandler(findSvc)
	commandBus.Register(find.AnswerCommandType, findAnswerCommandHandler)

	// update handler
	updateSvc := update.NewService(producer, storage, log)
	updateAnswerCommandHandler := update.NewAnswerCommandHandler(updateSvc)
	commandBus.Register(update.AnswerCommandType, updateAnswerCommandHandler)

	ctx, srv := New(ctx, cfg, commandBus)

	consumer := kafka.NewConsumer(strings.Split(cfg.Kafka.Consumer, ","), cfg.Kafka.Topic)
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
