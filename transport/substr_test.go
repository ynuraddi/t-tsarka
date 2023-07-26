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

func TestFindSubstring(t *testing.T) {
	testCases := []struct {
		name          string
		input         func() io.Reader
		buildStubs    func(service *mock_service.MockISubstrService)
		checkResponce func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			input: func() io.Reader {
				return strings.NewReader(`{"input":"abc"}`)
			},
			buildStubs: func(service *mock_service.MockISubstrService) {
				service.EXPECT().FindSubstr(gomock.Any()).Times(1).Return("abc")
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid JSON",
			input: func() io.Reader {
				return strings.NewReader(`{input}`)
			},
			buildStubs: func(service *mock_service.MockISubstrService) {
				service.EXPECT().FindSubstr(gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			input: func() io.Reader {
				return strings.NewReader(`{"input":"!!!!"}`)
			},
			buildStubs: func(service *mock_service.MockISubstrService) {
				service.EXPECT().FindSubstr(gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			substrService := mock_service.NewMockISubstrService(ctrl)
			test.buildStubs(substrService)

			service := &service.Manager{Substr: substrService}

			req := httptest.NewRequest(http.MethodPost, "/rest/substr/find", test.input())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			server := testServer(service)
			c := server.router.NewContext(req, rec)
			server.findSubstring(c)

			test.checkResponce(t, rec)
		})
	}
}
