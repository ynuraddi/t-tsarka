package service

import (
	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/ilogger"
	"github.com/ynuraddi/t-tsarka/service/substr"
)

type ISubstrService interface {
	FindSubstr(s string) string
}

type Manager struct {
	Substr ISubstrService
}

func New(config *config.Config, logger ilogger.ILogger) *Manager {
	substrService := substr.NewSubstrService(logger)

	manager := &Manager{
		Substr: substrService,
	}

	return manager
}
