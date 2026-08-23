package model

import (
	"context"
	"database/sql"
)

type Checkpoint struct {
	ChainID       int64
	Factory       []byte
	NextBlock     uint64
	LastBlockHash []byte
}

func (m *Model) GetCheckpoint(ctx context.Context, chainID int64, factory []byte) (Checkpoint, error) {
	var checkpoint Checkpoint
	err := m.db.QueryRowContext(ctx, `
		SELECT chain_id, factory, next_block, last_block_hash
		FROM indexer_checkpoints WHERE chain_id=$1 AND factory=$2`, chainID, factory,
	).Scan(&checkpoint.ChainID, &checkpoint.Factory, &checkpoint.NextBlock, &checkpoint.LastBlockHash)
	return checkpoint, err
}

func (m *Model) SaveCheckpoint(ctx context.Context, checkpoint Checkpoint) error {
	_, err := m.db.ExecContext(ctx, checkpointUpsertSQL,
		checkpoint.ChainID, checkpoint.Factory, checkpoint.NextBlock, nullableBytes(checkpoint.LastBlockHash),
	)
	return err
}

const checkpointUpsertSQL = `
	INSERT INTO indexer_checkpoints (chain_id, factory, next_block, last_block_hash)
	VALUES ($1,$2,$3,$4)
	ON CONFLICT (chain_id, factory) DO UPDATE SET
		next_block=EXCLUDED.next_block, last_block_hash=EXCLUDED.last_block_hash, updated_at=now()`

func saveCheckpointTx(ctx context.Context, tx *sql.Tx, checkpoint Checkpoint) error {
	_, err := tx.ExecContext(ctx, checkpointUpsertSQL,
		checkpoint.ChainID, checkpoint.Factory, checkpoint.NextBlock, nullableBytes(checkpoint.LastBlockHash),
	)
	return err
}

func (m *Model) DeleteLaunchEventsFromBlock(ctx context.Context, chainID int64, factory []byte, fromBlock uint64) error {
	return m.transact(ctx, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, `
			DELETE FROM token_launches
			WHERE chain_id=$1 AND factory=$2 AND block_number >= $3
			RETURNING tx_hash`, chainID, factory, fromBlock)
		if err != nil {
			return err
		}
		var hashes [][]byte
		for rows.Next() {
			var hash []byte
			if err := rows.Scan(&hash); err != nil {
				rows.Close()
				return err
			}
			hashes = append(hashes, hash)
		}
		if err := rows.Close(); err != nil {
			return err
		}
		for _, hash := range hashes {
			if _, err := tx.ExecContext(ctx, `
				UPDATE chain_transactions SET status='orphaned', block_number=NULL, block_hash=NULL, updated_at=now()
				WHERE chain_id=$1 AND tx_hash=$2`, chainID, hash); err != nil {
				return err
			}
		}
		return saveCheckpointTx(ctx, tx, Checkpoint{ChainID: chainID, Factory: factory, NextBlock: fromBlock})
	})
}
