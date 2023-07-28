package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ynuraddi/t-tsarka/service"
	mock_service "github.com/ynuraddi/t-tsarka/service/mock"
)

func TestCounterAdd(t *testing.T) {
	testCases := []struct {
		name          string
		param         string
		buildStubs    func(service *mock_service.MockICounterService)
		checkResponce func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: "1",
			buildStubs: func(service *mock_service.MockICounterService) {
				service.EXPECT().Add(gomock.Any()).Times(1).Return(nil)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "NaN Param",
			param: "abc",
			buildStubs: func(service *mock_service.MockICounterService) {
				service.EXPECT().Add(gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Internal",
			param: "1",
			buildStubs: func(service *mock_service.MockICounterService) {
				service.EXPECT().Add(gomock.Any()).Times(1).Return(errors.New("unexpected"))
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			counterService := mock_service.NewMockICounterService(ctrl)
			test.buildStubs(counterService)

			service := &service.Manager{Counter: counterService}

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()

			server := testServer(service)
			c := server.router.NewContext(req, rec)
			c.SetPath("/rest/counter/add/:i")
			c.SetParamNames("i")
			c.SetParamValues(test.param)

			server.counterAdd(c)

			test.checkResponce(t, rec)
		})
	}
}

func TestCounterSub(t *testing.T) {
	testCases := []struct {
		name          string
		param         string
		buildStubs    func(service *mock_service.MockICounterService)
		checkResponce func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: "1",
			buildStubs: func(service *mock_service.MockICounterService) {
				service.EXPECT().Sub(gomock.Any()).Times(1).Return(nil)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "NaN Param",
			param: "abc",
			buildStubs: func(service *mock_service.MockICounterService) {
				service.EXPECT().Sub(gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Internal",
			param: "1",
			buildStubs: func(service *mock_service.MockICounterService) {
				service.EXPECT().Sub(gomock.Any()).Times(1).Return(errors.New("unexpected"))
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			counterService := mock_service.NewMockICounterService(ctrl)
			test.buildStubs(counterService)

			service := &service.Manager{Counter: counterService}

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()

			server := testServer(service)
			c := server.router.NewContext(req, rec)
			c.SetPath("/rest/counter/add/:i")
			c.SetParamNames("i")
			c.SetParamValues(test.param)

			server.counterSub(c)

			test.checkResponce(t, rec)
		})
	}
}

func TestCounterGet(t *testing.T) {
	testCases := []struct {
		name          string
		buildStubs    func(service *mock_service.MockICounterService)
		checkResponce func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(service *mock_service.MockICounterService) {
				service.EXPECT().Get().Times(1).Return(int64(1), nil)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var get counterGetResponce
				err := json.NewDecoder(recorder.Body).Decode(&get)
				require.NoError(t, err)
				require.Equal(t, int64(1), get.Value)
			},
		},
		{
			name: "Internal",
			buildStubs: func(service *mock_service.MockICounterService) {
				service.EXPECT().Get().Times(1).Return(int64(0), errors.New("internal"))
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			counterService := mock_service.NewMockICounterService(ctrl)
			test.buildStubs(counterService)

			service := &service.Manager{Counter: counterService}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			server := testServer(service)
			c := server.router.NewContext(req, rec)

			server.counterGet(c)
			test.checkResponce(t, rec)
		})
	}
}
