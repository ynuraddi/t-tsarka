package email

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ynuraddi/t-tsarka/pkg/logger"
)

func TestEmailService(t *testing.T) {
	logger := logger.NewLogger(os.Stderr, logger.LvlTest, nil)
	service := NewEmailService(logger)

	testCases := []struct {
		input string
		want  []string
	}{
		{
			input: "Email:def@gmail.com",
			want:  []string{"def@gmail.com"},
		},
		{
			input: " Email:def@gmail.com",
			want:  []string{"def@gmail.com"},
		},
		{
			input: "Email:def@gmail.comEmail:def2@gmail.com",
			want:  nil,
		},
		{
			input: `Email:@gamil.com
			
			Email:def@gmail.com
			   Email:boba@pux.box`,
			want: []string{"def@gmail.com", "boba@pux.box"},
		},
	}

	for _, test := range testCases {
		t.Run("", func(t *testing.T) {
			get := service.Check(test.input)
			require.Equal(t, test.want, get)
		})
	}
}
