package service

import (
	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/ilogger"
)

type ISubstrService interface {
	FindSubstr(s string) string
}

type Manager struct {
	Substr ISubstrService
}

func New(config config.Config, logger ilogger.ILogger) *Manager {
	manager := &Manager{}

	return manager
}
