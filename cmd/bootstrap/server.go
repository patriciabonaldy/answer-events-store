package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	handler2 "github.com/patriciabonaldy/bequest_challenge/cmd/bootstrap/handler"
	"github.com/patriciabonaldy/bequest_challenge/cmd/docs"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-gonic/gin"
	"github.com/patriciabonaldy/bequest_challenge/internal/config"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine
	handler  handler2.AnswerHandler

	shutdownTimeout time.Duration
}

func New(ctx context.Context, config *config.Config, handler handler2.AnswerHandler) (context.Context, Server) {
	srv := Server{
		engine:   gin.New(),
		httpAddr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		handler:  handler,

		shutdownTimeout: time.Duration(config.ShutdownTimeout) + time.Second,
	}

	srv.registerRoutes()
	return serverContext(ctx), srv
}

// Middleware is a gin.HandlerFunc that set CORS
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func (s *Server) registerRoutes() {
	s.engine.Use(Middleware())
	s.engine.GET("/health", handler2.CheckHandler())
	answer := s.engine.Group("/answers")
	{
		answer.GET("/:id", s.handler.GetAnswer())
		answer.GET("/:id/history", s.handler.GetHistory())
		answer.POST("", s.handler.Create())
		answer.PUT("/:id", s.handler.Update())
		answer.DELETE("/:id", s.handler.Delete())
	}

	docs.SwaggerInfo.Title = "Swagger Documentation API"
	docs.SwaggerInfo.Description = "Event API Documentation."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "0.0.0.0:8080"
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
