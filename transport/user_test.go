package transport

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/ynuraddi/t-tsarka/model"
	"github.com/ynuraddi/t-tsarka/service"
	mock_service "github.com/ynuraddi/t-tsarka/service/mock"
)

func TestUserCreate(t *testing.T) {
	testCases := []struct {
		name  string
		input userCreateRequest

		buildStubs    func(mock *mock_service.MockIUserService)
		checkResponce func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			input: userCreateRequest{
				FirstName: "first",
				LastName:  "last",
			},
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(int64(1), nil)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				var resp userCreateResponce
				require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp))

				require.Equal(t, int64(1), resp.ID)
			},
		},
		{
			name: "BadRequest noFirstName",
			input: userCreateRequest{
				FirstName: "",
				LastName:  "last",
			},
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest noLastName",
			input: userCreateRequest{
				FirstName: "first",
				LastName:  "",
			},
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "BadRequest noContent",
			input: userCreateRequest{},
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal",
			input: userCreateRequest{
				FirstName: "first",
				LastName:  "last",
			},
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(int64(0), errors.New("internal"))
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

			userService := mock_service.NewMockIUserService(ctrl)
			test.buildStubs(userService)

			service := &service.Manager{User: userService}

			bodyJSON, err := json.Marshal(test.input)
			require.NoError(t, err)

			body := bytes.NewReader(bodyJSON)

			req := httptest.NewRequest(http.MethodPost, "/rest/user", body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			server := testServer(service)
			c := server.router.NewContext(req, rec)

			server.userCreate(c)

			test.checkResponce(t, rec)
		})
	}

	t.Run("unprocessible json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userService := mock_service.NewMockIUserService(ctrl)
		userService.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		service := &service.Manager{User: userService}

		bodyJSON, err := json.Marshal("")
		require.NoError(t, err)

		body := bytes.NewReader(bodyJSON)

		req := httptest.NewRequest(http.MethodPost, "/rest/user", body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		server := testServer(service)
		c := server.router.NewContext(req, rec)

		server.userCreate(c)

		require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})
}

func TestUserGet(t *testing.T) {
	testCases := []struct {
		name          string
		param         string
		buildStubs    func(mock *mock_service.MockIUserService)
		checkResponce func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: "1",
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(model.User{}, nil)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var resp model.User
				require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp))

				require.Equal(t, model.User{}, resp)
			},
		},
		{
			name:  "BadRequest incorrectParam",
			param: "abc",
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Get(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Internal",
			param: "1",
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(model.User{}, errors.New("internal"))
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "NotFound",
			param: "1",
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(model.User{}, sql.ErrNoRows)
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

			userService := mock_service.NewMockIUserService(ctrl)
			test.buildStubs(userService)

			service := &service.Manager{User: userService}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			server := testServer(service)
			c := server.router.NewContext(req, rec)
			c.SetPath("/rest/user/:id")
			c.SetParamNames("id")
			c.SetParamValues(test.param)

			server.userGet(c)

			test.checkResponce(t, rec)
		})
	}
}

func TestUserUpdate(t *testing.T) {
	testCases := []struct {
		name          string
		param         string
		input         string
		buildStubs    func(mock *mock_service.MockIUserService)
		checkResponce func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: "1",
			input: `{"first_name":"first", "last_name":"last"}`,
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(model.User{}, nil)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var resp model.User
				require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp))

				require.Equal(t, model.User{}, resp)
			},
		},
		{
			name:  "BadRequest incorrectParam",
			param: "abc",
			input: `{"first_name":"first", "last_name":"last"}`,
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "UnprocessibleJSON incorrectJSON",
			param: "1",
			input: `{"first}`,
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name:  "BadRequest failedValidate",
			param: "1",
			input: `{"first_name":"first", "last_name":""}`,
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Internal",
			param: "1",
			input: `{"first_name":"first", "last_name":"last"}`,
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(model.User{}, errors.New("internal"))
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "NotFound",
			param: "1",
			input: `{"first_name":"first", "last_name":"last"}`,
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Update(gomock.Any(), gomock.Any()).Times(1).Return(model.User{}, sql.ErrNoRows)
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

			userService := mock_service.NewMockIUserService(ctrl)
			test.buildStubs(userService)

			service := &service.Manager{User: userService}

			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(test.input))
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			server := testServer(service)
			c := server.router.NewContext(req, rec)
			c.SetPath("/rest/user/:id")
			c.SetParamNames("id")
			c.SetParamValues(test.param)

			server.userUpdate(c)

			test.checkResponce(t, rec)
		})
	}
}

func TestUserDelete(t *testing.T) {
	testCases := []struct {
		name          string
		param         string
		buildStubs    func(mock *mock_service.MockIUserService)
		checkResponce func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: "1",
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "BadRequest incorrectParam",
			param: "abc",
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Internal",
			param: "1",
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(1).Return(errors.New("internal"))
			},
			checkResponce: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "NotFound",
			param: "1",
			buildStubs: func(mock *mock_service.MockIUserService) {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(1).Return(sql.ErrNoRows)
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

			userService := mock_service.NewMockIUserService(ctrl)
			test.buildStubs(userService)

			service := &service.Manager{User: userService}

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			rec := httptest.NewRecorder()

			server := testServer(service)
			c := server.router.NewContext(req, rec)
			c.SetPath("/rest/user/:id")
			c.SetParamNames("id")
			c.SetParamValues(test.param)

			server.userDelete(c)

			test.checkResponce(t, rec)
		})
	}
}
