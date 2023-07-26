package transport

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/ilogger"
)

type Server struct {
	config config.Config

	logger ilogger.ILogger

	router *echo.Echo
}

func New(config config.Config, logger ilogger.ILogger) *Server {
	server := &Server{
		config: config,
		logger: logger,
	}

	return server
}

func (s *Server) Start(ctx context.Context) error {
	s.setupRouter()

	go func() {
		if err := s.router.Start(s.config.HttpHost + ":" + s.config.HttpPort); err != nil {
			s.logger.Fatal("failed start server", err)
		}
	}()
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := s.router.Shutdown(ctxShutDown); err != nil {
		s.logger.Fatal("failed shutdown server", err)
	}

	s.logger.Info("Sercer gracefully shutdown")
	return nil
}
