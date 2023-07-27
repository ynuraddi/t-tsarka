package iin

import (
	"regexp"

	"github.com/ynuraddi/t-tsarka/ilogger"
)

type iinService struct {
	logger ilogger.ILogger
}

func NewiinService(logger ilogger.ILogger) *iinService {
	return &iinService{
		logger: logger,
	}
}

const (
	IinRX = `\b\d{12}\b`
)

func (s *iinService) Check(str string) (result []string) {
	reg := regexp.MustCompile(IinRX)
	return reg.FindAllString(str, -1)
}
