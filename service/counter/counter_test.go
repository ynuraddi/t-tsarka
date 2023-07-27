package counter

import (
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ynuraddi/t-tsarka/config"
	"github.com/ynuraddi/t-tsarka/pkg/logger"
)

const counterKey = "counter"

func TestCounterService_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewLogger(nil, logger.LvlTest, nil)
	client, mock := redismock.NewClientMock()

	service := NewCounterService(&config.Config{
		RedisCounterKey: counterKey,
	}, logger, client)

	testCases := []struct {
		name      string
		input     int64
		buildStub func(i int64)
		want      error
	}{
		{
			name:  "OK",
			input: 1,
			buildStub: func(i int64) {
				mock.ExpectIncrBy(counterKey, i).SetVal(0)
			},
			want: nil,
		},
		{
			name:  "Internal",
			input: 1,
			buildStub: func(i int64) {
				mock.ExpectIncrBy(counterKey, i).SetErr(errors.New("internal"))
			},
			want: errors.New("internal"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.buildStub(test.input)
			err := service.Add(test.input)
			require.Equal(t, test.want, err)
		})
	}
}

func TestCounterService_Sub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewLogger(nil, logger.LvlTest, nil)
	client, mock := redismock.NewClientMock()

	service := NewCounterService(&config.Config{
		RedisCounterKey: counterKey,
	}, logger, client)

	testCases := []struct {
		name      string
		input     int64
		want      error
		buildStub func(i int64, err error)
	}{
		{
			name:  "OK",
			input: 1,
			want:  nil,
			buildStub: func(i int64, err error) {
				mock.ExpectDecrBy(counterKey, i).SetVal(0)
			},
		},
		{
			name:  "Internal",
			input: 1,
			want:  errors.New("internal"),
			buildStub: func(i int64, err error) {
				mock.ExpectDecrBy(counterKey, i).SetErr(errors.New("internal"))
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.buildStub(test.input, test.want)
			err := service.Sub(test.input)
			require.Equal(t, test.want, err)
		})
	}
}

func TestCounterService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewLogger(nil, logger.LvlTest, nil)
	client, mock := redismock.NewClientMock()

	service := NewCounterService(&config.Config{
		RedisCounterKey: counterKey,
	}, logger, client)

	testCases := []struct {
		name      string
		want      int
		wantErr   bool
		buildStub func()
	}{
		{
			name:    "OK",
			want:    1,
			wantErr: false,
			buildStub: func() {
				mock.ExpectGet(counterKey).SetVal("1")
			},
		},
		{
			name:    "Failed Get",
			want:    0,
			wantErr: true,
			buildStub: func() {
				mock.ExpectGet(counterKey).SetErr(errors.New("internal"))
			},
		},
		{
			name:    "Failed Convert",
			want:    0,
			wantErr: true,
			buildStub: func() {
				mock.ExpectGet(counterKey).SetVal("aboba")
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.buildStub()

			get, err := service.Get()
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.want, get)
		})
	}
}
