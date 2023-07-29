package transport

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/ilogger"
	"github.com/ynuraddi/t-tsarka/service"
)

type Server struct {
	config *config.Config
	logger ilogger.ILogger

	router  *echo.Echo
	service *service.Manager
}

func New(config *config.Config, logger ilogger.ILogger, service *service.Manager) *Server {
	server := &Server{
		config: config,
		logger: logger,

		service: service,
	}

	return server
}

func (s *Server) Start(ctx context.Context) error {
	s.setupRouter()

	go func() {
		if err := s.router.Start(s.config.HttpHost + ":" + s.config.HttpPort); err != nil {
			s.logger.Fatal("server stoped", err)
		}
	}()
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	s.logger.Debug("gracefull shutdown is started")
	if err := s.router.Shutdown(ctxShutDown); err != nil {
		s.logger.Fatal("failed shutdown server", err)
	}

	s.logger.Info("Server gracefully shutdown")
	return nil
}
