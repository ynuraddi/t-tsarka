package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ynuraddi/t-tsarka/ilogger"
	"github.com/ynuraddi/t-tsarka/model"
)

type userRepository struct {
	logger ilogger.ILogger
	db     *sql.DB
}

func NewUserRepository(logger ilogger.ILogger, db *sql.DB) *userRepository {
	return &userRepository{
		logger: logger,
		db:     db,
	}
}

func (r *userRepository) Create(ctx context.Context, user model.User) (id int64, err error) {
	query := "insert into users (first_name, last_name) values ($1, $2) returning id"

	if err := r.db.QueryRowContext(ctx, query, user.FirstName, user.LastName).Scan(&id); err != nil {
		r.logger.Error("failed insert user", err)
		return 0, err
	}

	return id, nil
}

func (r *userRepository) Get(ctx context.Context, id int64) (dbuser model.User, err error) {
	query := "select * from users where id = $1"

	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dbuser.ID,
		&dbuser.FirstName,
		&dbuser.LastName,
	); err != nil {
		r.logger.Error("failed get user", err)
		return model.User{}, err
	}

	return dbuser, nil
}

// TODO add versioning or tx with locking
func (r *userRepository) Update(ctx context.Context, user model.User) (dbuser model.User, err error) {
	query := `update users set
	first_name = case when $1 = '' then first_name else $1 end,
	last_name = case when $2 = '' then last_name else $2 end
	where id = $3
	returning id, first_name, last_name`

	if err = r.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.ID).Scan(
		&dbuser.ID, &dbuser.FirstName, &dbuser.LastName,
	); err != nil {
		r.logger.Error("failed to update user", err)
		return model.User{}, err
	}

	return dbuser, nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `delete from users where id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete user", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows, delete user", err)
		return err
	}

	if rowsAffected == 0 {
		err := fmt.Errorf("no affected rows: %w", sql.ErrNoRows)
		r.logger.Error("failed to delete user", err)
		return err
	}

	return nil
}
