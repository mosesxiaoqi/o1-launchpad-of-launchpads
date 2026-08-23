package database

import (
	"context"
	"database/sql"
	"fmt"

	commonconfig "o1-launchpad/common/config"

	_ "github.com/lib/pq"
)

func Open(ctx context.Context, config commonconfig.Database) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DataSource)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect PostgreSQL: %w", err)
	}
	return db, nil
}
