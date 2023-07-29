package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ynuraddi/t-tsarka/config"

	_ "github.com/lib/pq"
)

func Open(config *config.Config) (*sql.DB, error) {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresUser,
		config.PostgresPass,
		config.PostgresDBName)

	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
