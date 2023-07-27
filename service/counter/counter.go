package counter

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/ilogger"
)

type counterService struct {
	key    string
	logger ilogger.ILogger
	client *redis.Client
}

func NewCounterService(config *config.Config, logger ilogger.ILogger, client *redis.Client) *counterService {
	return &counterService{
		key:    config.RedisCounterKey,
		logger: logger,
		client: client,
	}
}

func (s *counterService) Add(i int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.client.IncrBy(ctx, s.key, i).Result()
	if err != nil {
		s.logger.Error("failed to increment CounterService", err)
		return err
	}
	s.logger.Debug("success increment CounterService")
	return nil
}

func (s *counterService) Sub(i int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.client.DecrBy(ctx, s.key, i).Result()
	if err != nil {
		s.logger.Error("failed to dicrement CounterService", err)
		return err
	}
	s.logger.Debug("success dicrement CounterService")
	return nil
}

func (s *counterService) Get() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := s.client.Get(ctx, s.key).Result()
	if err != nil {
		s.logger.Error("failed to get value CounterService", err)
		return 0, err
	}
	num, err := strconv.Atoi(result)
	if err != nil {
		s.logger.Error("failed to convert value CounterService", err)
		return 0, err
	}
	return num, nil
}
