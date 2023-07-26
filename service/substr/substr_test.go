package substr

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ynuraddi/t-tsarka/pkg/logger"
)

func TestSubstrService(t *testing.T) {
	logger := logger.NewLogger(os.Stderr, logger.LvlTest, nil)
	service := NewSubstrService(logger)

	testCases := []struct {
		input     string
		wantOneOf []string
	}{
		{
			input:     "abc",
			wantOneOf: []string{"abc"},
		},
		{
			input:     "geekforgeek",
			wantOneOf: []string{"ekforg", "kforge"},
		},
		{
			input:     "abbbbbaaaaba",
			wantOneOf: []string{"ab", "ba"},
		},
		{
			input:     "aaaaboe",
			wantOneOf: []string{"aboe"},
		},
		{
			input:     "1234567890123",
			wantOneOf: []string{"1234567890", "4567890123"},
		},
	}

	for _, test := range testCases {
		t.Run("", func(t *testing.T) {
			get := service.FindSubstr(test.input)
			require.Contains(t, test.wantOneOf, get, "output: %s", get)
		})
	}
}
