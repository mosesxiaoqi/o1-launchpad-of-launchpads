package model

import (
	"context"
	"database/sql"
	"time"
)

type ChainTransaction struct {
	ChainID       int64
	TxHash        []byte
	Kind          string
	Status        string
	BlockNumber   uint64
	BlockHash     []byte
	FailureReason string
	UpdatedAt     time.Time
}

func (m *Model) UpsertTransaction(ctx context.Context, transaction ChainTransaction) error {
	_, err := m.db.ExecContext(ctx, transactionUpsertSQL,
		transaction.ChainID, transaction.TxHash, transaction.Kind, transaction.Status,
		nullableBlock(transaction.BlockNumber), nullableBytes(transaction.BlockHash), transaction.FailureReason,
	)
	return err
}

func (m *Model) GetTransaction(ctx context.Context, chainID int64, txHash []byte) (ChainTransaction, error) {
	var transaction ChainTransaction
	var blockNumber sql.NullInt64
	err := m.db.QueryRowContext(ctx, `
		SELECT chain_id, tx_hash, kind, status, block_number, block_hash, failure_reason, updated_at
		FROM chain_transactions WHERE chain_id=$1 AND tx_hash=$2`, chainID, txHash,
	).Scan(
		&transaction.ChainID, &transaction.TxHash, &transaction.Kind, &transaction.Status,
		&blockNumber, &transaction.BlockHash, &transaction.FailureReason, &transaction.UpdatedAt,
	)
	if blockNumber.Valid {
		transaction.BlockNumber = uint64(blockNumber.Int64)
	}
	return transaction, err
}

const transactionUpsertSQL = `
	INSERT INTO chain_transactions (
		chain_id, tx_hash, kind, status, block_number, block_hash, failure_reason
	) VALUES ($1,$2,$3,$4,$5,$6,$7)
	ON CONFLICT (chain_id, tx_hash) DO UPDATE SET
		kind=EXCLUDED.kind, status=EXCLUDED.status, block_number=EXCLUDED.block_number,
		block_hash=EXCLUDED.block_hash, failure_reason=EXCLUDED.failure_reason, updated_at=now()`

func upsertTransactionTx(ctx context.Context, tx *sql.Tx, transaction ChainTransaction) error {
	_, err := tx.ExecContext(ctx, transactionUpsertSQL,
		transaction.ChainID, transaction.TxHash, transaction.Kind, transaction.Status,
		nullableBlock(transaction.BlockNumber), nullableBytes(transaction.BlockHash), transaction.FailureReason,
	)
	return err
}

func nullableBlock(block uint64) any {
	if block == 0 {
		return nil
	}
	return block
}

func nullableBytes(value []byte) any {
	if len(value) == 0 {
		return nil
	}
	return value
}
