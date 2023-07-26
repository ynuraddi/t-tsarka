package email

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ynuraddi/t-tsarka/ilogger"
)

type emailService struct {
	logger ilogger.ILogger
}

func NewEmailService(logger ilogger.ILogger) *emailService {
	return &emailService{
		logger: logger,
	}
}

const (
	EmailRX     = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	EmailPrefix = "Email:"
)

func (s *emailService) Check(str string) (result []string) {
	s.logger.Debug(fmt.Sprintf("start find email in \"%s\"", str))

	reg := regexp.MustCompile(EmailRX)

	datas := strings.Fields(str)
	for _, data := range datas {
		if strings.HasPrefix(data, EmailPrefix) {
			data = strings.TrimPrefix(data, EmailPrefix)
			if reg.MatchString(data) {
				result = append(result, data)
			}
		}
	}
	return
}
