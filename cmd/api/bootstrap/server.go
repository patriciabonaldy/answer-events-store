package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/patriciabonaldy/bequest_challenge/cmd/api/docs"

	"github.com/gin-gonic/gin"
	"github.com/patriciabonaldy/bequest_challenge/cmd/api/bootstrap/handler"
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
	srv.setCors()
	return serverContext(ctx), srv
}

func (s *Server) setCors() {
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))
}
func (s *Server) registerRoutes() {
	s.engine.GET("/health", handler.CheckHandler())
	answer := s.engine.Group("/answer")
	{
		answer.GET("/:id", s.handler.Get())
		answer.GET("/:id/history", s.handler.GetHistory())
		answer.POST("", s.handler.Create())
		answer.PUT("/:id", s.handler.Update())
		answer.DELETE("/:id", s.handler.Delete())
	}

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "Event API Documentation."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// use ginSwagger middleware to serve the API docs
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
