package transport

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/ynuraddi/t-tsarka/service"
	mock_service "github.com/ynuraddi/t-tsarka/service/mock"
)

func TestEmailCheck(t *testing.T) {
	testCases := []struct {
		name          string
		input         func() io.Reader
		buildStubs    func(service *mock_service.MockIEmailService)
		checkResponce func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			input: func() io.Reader {
				return strings.NewReader(`{"input":"def@gmail.com"}`)
			},
			buildStubs: func(service *mock_service.MockIEmailService) {
				service.EXPECT().Check(gomock.Any()).Times(1).Return([]string{"def@gmail.com"})
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid JSON",
			input: func() io.Reader {
				return strings.NewReader(`{"input"}`)
			},
			buildStubs: func(service *mock_service.MockIEmailService) {
				service.EXPECT().Check(gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			input: func() io.Reader {
				return strings.NewReader(`{"input":""}`)
			},
			buildStubs: func(service *mock_service.MockIEmailService) {
				service.EXPECT().Check(gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found",
			input: func() io.Reader {
				return strings.NewReader(`{"input":"def def"}`)
			},
			buildStubs: func(service *mock_service.MockIEmailService) {
				service.EXPECT().Check(gomock.Any()).Times(1).Return(nil)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			emailService := mock_service.NewMockIEmailService(ctrl)
			test.buildStubs(emailService)

			service := &service.Manager{Email: emailService}

			req := httptest.NewRequest(http.MethodPost, "/rest/email/check", test.input())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			server := testServer(service)
			c := server.router.NewContext(req, rec)
			server.emailCheck(c)

			test.checkResponce(t, rec)
		})
	}
}
