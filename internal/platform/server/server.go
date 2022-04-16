package server

import (
	"context"
	"fmt"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/server/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/patriciabonaldy/bequest_challenge/internal/config"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine
	handler  handler.AnswerHandler

	shutdownTimeout time.Duration
}

func New(ctx context.Context, config *config.Config, handler handler.AnswerHandler) (context.Context, Server) {
	srv := Server{
		engine:   gin.New(),
		httpAddr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		handler:  handler,

		shutdownTimeout: time.Duration(config.ShutdownTimeout) + time.Second,
	}

	srv.registerRoutes()
	return serverContext(ctx), srv
}

func (s *Server) registerRoutes() {
	s.engine.GET("/health", handler.CheckHandler())
	s.engine.POST("/answer", s.handler.Create())
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
