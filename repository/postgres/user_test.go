package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/ynuraddi/t-tsarka/model"
	"github.com/ynuraddi/t-tsarka/pkg/logger"
)

func TestUserRepoCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := logger.NewLogger(nil, logger.LvlTest, nil)
	repo := NewUserRepository(logger, db)

	testCases := []struct {
		name        string
		firstName   string
		lastName    string
		buildStub   func()
		checkResult func(get int64, err error)
	}{
		{
			name:      "OK",
			firstName: "aboba1",
			lastName:  "aboba2",
			buildStub: func() {
				mock.ExpectExec("insert into table users").WithArgs("aboba1", "aboba2").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			checkResult: func(get int64, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(1), get)
			},
		},
		{
			name:      "Exec error",
			firstName: "aboba1",
			lastName:  "aboba2",
			buildStub: func() {
				mock.ExpectExec("insert into table users").WithArgs("aboba1", "aboba2").WillReturnError(errors.New("internal"))
			},
			checkResult: func(get int64, err error) {
				require.Error(t, err)
				require.Equal(t, int64(0), get)
			},
		},
		{
			name:      "LastInsertId error",
			firstName: "aboba1",
			lastName:  "aboba2",
			buildStub: func() {
				mock.ExpectExec("insert into table users").WithArgs("aboba1", "aboba2").WillReturnResult(sqlmock.NewErrorResult(errors.New("internal")))
			},
			checkResult: func(get int64, err error) {
				require.Error(t, err)
				require.Equal(t, int64(0), get)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.buildStub()
			get, err := repo.Create(context.Background(), model.User{
				FirstName: test.firstName,
				LastName:  test.lastName,
			})
			test.checkResult(get, err)
		})
	}
}

func TestUserRepoGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := logger.NewLogger(nil, logger.LvlTest, nil)
	repo := NewUserRepository(logger, db)

	reqexpr := `select\s+\*\s+from\s+users\s+where\s+id\s+=\s+\$1`
	columns := []string{"id", "first_name", "last_name"}

	testCases := []struct {
		name        string
		input       int64
		buildStub   func()
		checkResult func(get model.User, err error)
	}{
		{
			name:  "OK",
			input: 1,
			buildStub: func() {
				mock.ExpectQuery(reqexpr).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, "first", "last"))
			},
			checkResult: func(get model.User, err error) {
				require.NoError(t, err)
				require.Equal(t, model.User{
					ID:        1,
					FirstName: "first",
					LastName:  "last",
				}, get)
			},
		},
		{
			name:  "Query Error",
			input: 1,
			buildStub: func() {
				mock.ExpectQuery(reqexpr).WithArgs(1).WillReturnError(errors.New("internal"))
			},
			checkResult: func(get model.User, err error) {
				require.Error(t, err)
				require.Equal(t, model.User{}, get)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.buildStub()
			user, err := repo.Get(context.Background(), test.input)
			test.checkResult(user, err)
		})
	}
}

func TestUserRepoUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := logger.NewLogger(nil, logger.LvlTest, nil)
	repo := NewUserRepository(logger, db)

	reqexpr := `update\s+users`
	columns := []string{"id", "first_name", "last_name"}

	testCases := []struct {
		name        string
		input       model.User
		buildStub   func()
		checkResult func(get model.User, err error)
	}{
		{
			name: "OK",
			input: model.User{
				ID:        1,
				FirstName: "first",
				LastName:  "last",
			},
			buildStub: func() {
				mock.ExpectQuery(reqexpr).WithArgs("first", "last", 1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, "first", "last"))
			},
			checkResult: func(get model.User, err error) {
				require.NoError(t, err)
				require.Equal(t, model.User{
					ID:        1,
					FirstName: "first",
					LastName:  "last",
				}, get)
			},
		},
		{
			name: "Query Error",
			input: model.User{
				ID:        1,
				FirstName: "",
				LastName:  "last",
			},
			buildStub: func() {
				mock.ExpectQuery(reqexpr).WithArgs("", "last", 1).WillReturnError(errors.New("internal"))
			},
			checkResult: func(get model.User, err error) {
				require.Error(t, err)
				require.Equal(t, model.User{}, get)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.buildStub()
			user, err := repo.Update(context.Background(), test.input)
			test.checkResult(user, err)
		})
	}
}

func TestUserRepoDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := logger.NewLogger(nil, logger.LvlTest, nil)
	repo := NewUserRepository(logger, db)

	reqexpr := `delete\s+from\s+users`

	testCases := []struct {
		name        string
		input       int64
		buildStub   func()
		checkResult func(err error)
	}{
		{
			name:  "OK",
			input: 1,
			buildStub: func() {
				mock.ExpectExec(reqexpr).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			checkResult: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name:  "Exec Error",
			input: 1,
			buildStub: func() {
				mock.ExpectExec(reqexpr).WithArgs(1).WillReturnError(errors.New("internal"))
			},
			checkResult: func(err error) {
				require.Error(t, err)
			},
		},
		{
			name:  "rowsAffected Error",
			input: 1,
			buildStub: func() {
				mock.ExpectExec(reqexpr).WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(errors.New("result error")))
			},
			checkResult: func(err error) {
				require.Error(t, err)
			},
		},
		{
			name:  "no affected rows",
			input: 1,
			buildStub: func() {
				mock.ExpectExec(reqexpr).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 0))
			},
			checkResult: func(err error) {
				require.Error(t, err)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.buildStub()
			err := repo.Delete(context.Background(), test.input)
			test.checkResult(err)
		})
	}
}
