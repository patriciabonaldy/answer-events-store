package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/patriciabonaldy/bequest_challenge/internal/config"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/server/handler/health"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine

	shutdownTimeout time.Duration
}

func New(ctx context.Context, config *config.Config) (context.Context, Server) {
	srv := Server{
		engine:   gin.New(),
		httpAddr: fmt.Sprintf("%s:%d", config.Host, config.Port),

		shutdownTimeout: config.ShutdownTimeout,
	}

	srv.registerRoutes()
	return serverContext(ctx), srv
}

func (s *Server) registerRoutes() {
	s.engine.GET("/health", health.CheckHandler())
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Server running on", s.httpAddr)

	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}
