package model

import (
	"context"
	"database/sql"
	"errors"
)

var ErrOwnershipConflict = errors.New("immutable launch ownership conflict")

type Model struct {
	db *sql.DB
}

func New(db *sql.DB) *Model {
	return &Model{db: db}
}

func (m *Model) transact(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
