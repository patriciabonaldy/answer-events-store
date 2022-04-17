package bootstrap

import (
	"context"
	"github.com/patriciabonaldy/bequest_challenge/cmd/api/bootstrap/handler"
	"github.com/patriciabonaldy/bequest_challenge/internal/business"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/storage/mongo"

	"github.com/patriciabonaldy/bequest_challenge/internal/config"
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

	svc := business.NewService(storage, log)
	handler := handler.New(svc, log)
	ctx, srv := New(ctx, cfg, handler)

	return srv.Run(ctx)
}
