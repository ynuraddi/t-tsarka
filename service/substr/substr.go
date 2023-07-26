package substr

import (
	"fmt"

	"github.com/ynuraddi/t-tsarka/ilogger"
)

type substrService struct {
	logger ilogger.ILogger
}

func NewSubstrService(logger ilogger.ILogger) *substrService {
	return &substrService{
		logger: logger,
	}
}

func (s *substrService) FindSubstr(str string) string {
	s.logger.Debug(fmt.Sprintf("start find substr in \"%s\"", str))

	var result string

	// 26 alpha characters, with upper case * 2, with numb + 10
	hash := make(map[byte]int, 26*2+10)
	var countOfDuplicate int

	var currentMaximum int

	for i, j := 0, 0; j < len(str); j++ {
		hash[str[j]] += 1

		if hash[str[j]] > 1 {
			countOfDuplicate += 1
		}

		if countOfDuplicate > 0 {
			if hash[str[i]] > 1 {
				countOfDuplicate -= 1
			}
			hash[str[i]] -= 1
			i += 1
		}

		if countOfDuplicate == 0 && j-i >= currentMaximum {
			result = str[i : j+1]
		}
	}

	return result
}
