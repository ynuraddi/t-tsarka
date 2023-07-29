package repository

import (
	"context"

	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/ilogger"
	"github.com/ynuraddi/t-tsarka/model"
	"github.com/ynuraddi/t-tsarka/repository/postgres"
)

type IUserRepository interface {
	Create(ctx context.Context, user model.User) (id int64, err error)
	Get(ctx context.Context, id int64) (user model.User, err error)
	Update(ctx context.Context, user model.User) (dbuser model.User, err error)
	Delete(ctx context.Context, id int64) error
}

type Manager struct {
	User IUserRepository
}

func New(config *config.Config, logger ilogger.ILogger) (*Manager, error) {
	db, err := postgres.Open(config)
	if err != nil {
		logger.Error("failed init postgres DB", err)
		return nil, err
	}
	logger.Debug("postgres success inited")

	userRepostiory := postgres.NewUserRepository(logger, db)

	return &Manager{
		User: userRepostiory,
	}, nil
}
