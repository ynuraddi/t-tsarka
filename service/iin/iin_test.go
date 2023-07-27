package iin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ynuraddi/t-tsarka/pkg/logger"
)

func TestIINService(t *testing.T) {
	logger := logger.NewLogger(os.Stderr, logger.LvlTest, nil)
	service := NewiinService(logger)

	testCases := []struct {
		input    string
		expected []string
	}{
		{input: "ИИН Казахстана: 920609301111, 860701402016", expected: []string{"920609301111", "860701402016"}},
		{input: "Some text without IIN", expected: nil},
		{input: "Invalid IIN format: 1234567890123", expected: nil},
		{input: "IINs separated by spaces: 123456789012 860701402016 920609301111", expected: []string{"123456789012", "860701402016", "920609301111"}},
	}

	for _, test := range testCases {
		t.Run("", func(t *testing.T) {
			get := service.Check(test.input)
			require.Equal(t, test.expected, get)
		})
	}
}
