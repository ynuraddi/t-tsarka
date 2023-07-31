package service

import (
	"context"

	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/ilogger"
	"github.com/ynuraddi/t-tsarka/model"
	"github.com/ynuraddi/t-tsarka/pkg/client/redis"
	"github.com/ynuraddi/t-tsarka/repository"
	"github.com/ynuraddi/t-tsarka/service/counter"
	"github.com/ynuraddi/t-tsarka/service/email"
	"github.com/ynuraddi/t-tsarka/service/iin"
	"github.com/ynuraddi/t-tsarka/service/substr"
	"github.com/ynuraddi/t-tsarka/service/user"
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

type IUserService interface {
	Create(ctx context.Context, user model.User) (id int64, err error)
	Get(ctx context.Context, id int64) (user model.User, err error)
	Update(ctx context.Context, user model.User) (dbuser model.User, err error)
	Delete(ctx context.Context, id int64) error
}

type Manager struct {
	Substr  ISubstrService
	Email   IEmailService
	IIN     IIINService
	Counter ICounterService
	User    IUserService
}

func New(config *config.Config, logger ilogger.ILogger, repo *repository.Manager) (*Manager, error) {
	redisClient, err := redis.NewClient(config)
	if err != nil {
		logger.Error("failed init redis client", err)
		return nil, err
	}

	userService := user.New(repo.User)

	substrService := substr.NewSubstrService(logger)
	emailService := email.NewEmailService(logger)
	iinService := iin.NewiinService(logger)

	counterService := counter.NewCounterService(config, logger, redisClient)

	manager := &Manager{
		Substr:  substrService,
		Email:   emailService,
		IIN:     iinService,
		Counter: counterService,
		User:    userService,
	}

	return manager, nil
}
