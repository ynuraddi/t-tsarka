package transport

import (
	"os"

	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/pkg/logger"
	"github.com/ynuraddi/t-tsarka/service"
)

func testServer(service *service.Manager) *Server {
	config := &config.Config{}
	logger := logger.NewLogger(os.Stderr, logger.LvlTest, nil)

	server := New(config, logger, service)
	server.setupRouter()

	return server
}
