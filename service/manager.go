package service

import (
	"github.com/redis/go-redis/v9"
	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/ilogger"
	"github.com/ynuraddi/t-tsarka/service/counter"
	"github.com/ynuraddi/t-tsarka/service/email"
	"github.com/ynuraddi/t-tsarka/service/iin"
	"github.com/ynuraddi/t-tsarka/service/substr"
)

type ISubstrService interface {
	FindSubstr(s string) string
}

type IEmailService interface {
	// return nil if found nothing
	Check(s string) []string
}

type IIINService interface {
	// return nil if found nothing
	Check(s string) []string
}

type ICounterService interface {
	Add(i int64) error
	Sub(i int64) error
	Get() (int64, error)
}

type Manager struct {
	Substr  ISubstrService
	Email   IEmailService
	IIN     IIINService
	Counter ICounterService
}

func New(config *config.Config, logger ilogger.ILogger, redisClient *redis.Client) *Manager {
	substrService := substr.NewSubstrService(logger)
	emailService := email.NewEmailService(logger)
	iinService := iin.NewiinService(logger)

	counterService := counter.NewCounterService(config, logger, redisClient)

	manager := &Manager{
		Substr:  substrService,
		Email:   emailService,
		IIN:     iinService,
		Counter: counterService,
	}

	return manager
}
