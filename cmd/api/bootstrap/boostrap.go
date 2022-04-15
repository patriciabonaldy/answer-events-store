package bootstrap

import (
	"context"

	"github.com/patriciabonaldy/bequest_challenge/internal/config"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/server"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	ctx, srv := server.New(context.Background(), cfg)
	return srv.Run(ctx)
}
