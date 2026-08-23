package model

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"o1-launchpad/common/pagination"

	"github.com/lib/pq"
)

type LaunchEvent struct {
	ChainID     int64
	LaunchpadID []byte
	Token       []byte
	PoolID      []byte
	Creator     []byte
	Quote       []byte
	Factory     []byte
	Supply      string
	TxHash      []byte
	BlockNumber uint64
	BlockHash   []byte
	LogIndex    uint32
	Status      string
	CreatedAt   time.Time
}

type TokenLaunch struct {
	ID int64
	LaunchEvent
	LaunchpadSlug string
}

func (m *Model) GetToken(ctx context.Context, chainID int64, token []byte) (TokenLaunch, error) {
	var item TokenLaunch
	err := m.db.QueryRowContext(ctx, `
		SELECT t.id, t.chain_id, p.launchpad_id, t.token, t.pool_id, t.creator, t.quote, t.factory,
		       t.supply::text, t.tx_hash, t.block_number, t.block_hash, t.log_index, t.status,
		       t.created_at, p.slug
		FROM token_launches t JOIN launchpads p ON p.id=t.launchpad_pk
		WHERE t.chain_id=$1 AND t.token=$2`, chainID, token,
	).Scan(
		&item.ID, &item.ChainID, &item.LaunchpadID, &item.Token, &item.PoolID, &item.Creator, &item.Quote,
		&item.Factory, &item.Supply, &item.TxHash, &item.BlockNumber, &item.BlockHash,
		&item.LogIndex, &item.Status, &item.CreatedAt, &item.LaunchpadSlug,
	)
	return item, err
}

func (m *Model) ListTokens(ctx context.Context, chainID int64, limit int) ([]TokenLaunch, error) {
	return m.ListTokensPage(ctx, chainID, limit, nil)
}

func (m *Model) ListTokensPage(ctx context.Context, chainID int64, limit int, cursor *pagination.Cursor) ([]TokenLaunch, error) {
	query := `
		SELECT t.id, t.chain_id, p.launchpad_id, t.token, t.pool_id, t.creator, t.quote, t.factory,
		       t.supply::text, t.tx_hash, t.block_number, t.block_hash, t.log_index, t.status,
		       t.created_at, p.slug
		FROM token_launches t JOIN launchpads p ON p.id=t.launchpad_pk
		WHERE t.chain_id=$1`
	args := []any{chainID}
	if cursor != nil {
		query += ` AND (t.created_at, t.id) < ($2, $3)`
		args = append(args, cursor.CreatedAt, cursor.ID)
	}
	args = append(args, limit)
	query += fmt.Sprintf(` ORDER BY t.created_at DESC, t.id DESC LIMIT $%d`, len(args))
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var launches []TokenLaunch
	for rows.Next() {
		var item TokenLaunch
		if err := rows.Scan(
			&item.ID, &item.ChainID, &item.LaunchpadID, &item.Token, &item.PoolID, &item.Creator, &item.Quote,
			&item.Factory, &item.Supply, &item.TxHash, &item.BlockNumber, &item.BlockHash,
			&item.LogIndex, &item.Status, &item.CreatedAt, &item.LaunchpadSlug,
		); err != nil {
			return nil, err
		}
		launches = append(launches, item)
	}
	return launches, rows.Err()
}

func (m *Model) ApplyLaunchEvent(ctx context.Context, event LaunchEvent, checkpoint Checkpoint) error {
	return m.transact(ctx, func(tx *sql.Tx) error {
		existing, err := getLaunchEventTx(ctx, tx, event.ChainID, event.Token)
		switch {
		case err == nil:
			if !sameLaunch(existing, event) {
				return ErrOwnershipConflict
			}
		case !errors.Is(err, sql.ErrNoRows):
			return err
		default:
			launchpadPK, err := launchpadPrimaryKey(ctx, tx, event.ChainID, event.LaunchpadID)
			if err != nil {
				return err
			}
			_, err = tx.ExecContext(ctx, `
				INSERT INTO token_launches (
					launchpad_pk, chain_id, token, pool_id, creator, quote, factory, supply,
					tx_hash, block_number, block_hash, log_index, status, created_at
				) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
				launchpadPK, event.ChainID, event.Token, event.PoolID, event.Creator, event.Quote,
				event.Factory, event.Supply, event.TxHash, event.BlockNumber, event.BlockHash,
				event.LogIndex, event.Status, event.CreatedAt,
			)
			if err != nil {
				var pqErr *pq.Error
				if errors.As(err, &pqErr) && pqErr.Code == "23505" {
					return ErrOwnershipConflict
				}
				return err
			}
		}

		if err := upsertTransactionTx(ctx, tx, ChainTransaction{
			ChainID: event.ChainID, TxHash: event.TxHash, Kind: "launch", Status: event.Status,
			BlockNumber: event.BlockNumber, BlockHash: event.BlockHash,
		}); err != nil {
			return err
		}
		return saveCheckpointTx(ctx, tx, checkpoint)
	})
}

func (m *Model) ListLaunchpadTokens(ctx context.Context, chainID int64, slug string, limit int) ([]TokenLaunch, error) {
	return m.ListLaunchpadTokensPage(ctx, chainID, slug, limit, nil)
}

func (m *Model) ListLaunchpadTokensPage(ctx context.Context, chainID int64, slug string, limit int, cursor *pagination.Cursor) ([]TokenLaunch, error) {
	query := `
		SELECT t.id, t.chain_id, p.launchpad_id, t.token, t.pool_id, t.creator, t.quote, t.factory,
		       t.supply::text, t.tx_hash, t.block_number, t.block_hash, t.log_index, t.status,
		       t.created_at, p.slug
		FROM token_launches t JOIN launchpads p ON p.id=t.launchpad_pk
		WHERE t.chain_id=$1 AND p.slug=$2`
	args := []any{chainID, slug}
	if cursor != nil {
		query += ` AND (t.created_at, t.id) < ($3, $4)`
		args = append(args, cursor.CreatedAt, cursor.ID)
	}
	args = append(args, limit)
	query += fmt.Sprintf(` ORDER BY t.created_at DESC, t.id DESC LIMIT $%d`, len(args))
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var launches []TokenLaunch
	for rows.Next() {
		var item TokenLaunch
		if err := rows.Scan(
			&item.ID, &item.ChainID, &item.LaunchpadID, &item.Token, &item.PoolID, &item.Creator, &item.Quote,
			&item.Factory, &item.Supply, &item.TxHash, &item.BlockNumber, &item.BlockHash,
			&item.LogIndex, &item.Status, &item.CreatedAt, &item.LaunchpadSlug,
		); err != nil {
			return nil, err
		}
		launches = append(launches, item)
	}
	return launches, rows.Err()
}

func getLaunchEventTx(ctx context.Context, tx *sql.Tx, chainID int64, token []byte) (LaunchEvent, error) {
	var event LaunchEvent
	err := tx.QueryRowContext(ctx, `
		SELECT t.chain_id, p.launchpad_id, t.token, t.pool_id, t.creator, t.quote, t.factory,
		       t.supply::text, t.tx_hash, t.block_number, t.block_hash, t.log_index, t.status, t.created_at
		FROM token_launches t JOIN launchpads p ON p.id=t.launchpad_pk
		WHERE t.chain_id=$1 AND t.token=$2`, chainID, token,
	).Scan(
		&event.ChainID, &event.LaunchpadID, &event.Token, &event.PoolID, &event.Creator, &event.Quote,
		&event.Factory, &event.Supply, &event.TxHash, &event.BlockNumber, &event.BlockHash,
		&event.LogIndex, &event.Status, &event.CreatedAt,
	)
	return event, err
}

func sameLaunch(a, b LaunchEvent) bool {
	return a.ChainID == b.ChainID && bytes.Equal(a.LaunchpadID, b.LaunchpadID) && bytes.Equal(a.Token, b.Token) &&
		bytes.Equal(a.PoolID, b.PoolID) && bytes.Equal(a.Creator, b.Creator) && bytes.Equal(a.Quote, b.Quote) &&
		bytes.Equal(a.Factory, b.Factory) && a.Supply == b.Supply && bytes.Equal(a.TxHash, b.TxHash) &&
		a.BlockNumber == b.BlockNumber && bytes.Equal(a.BlockHash, b.BlockHash) && a.LogIndex == b.LogIndex
}
