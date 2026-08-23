package model

import (
	"context"
	"database/sql"
	"time"
)

type AuthNonce struct {
	ID        string
	ChainID   int64
	Wallet    []byte
	NonceHash []byte
	Message   string
	ExpiresAt time.Time
}

func (m *Model) CreateAuthNonce(ctx context.Context, nonce AuthNonce) error {
	_, err := m.db.ExecContext(ctx, `
		INSERT INTO auth_nonces (id, chain_id, wallet, nonce_hash, message, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		nonce.ID, nonce.ChainID, nonce.Wallet, nonce.NonceHash, nonce.Message, nonce.ExpiresAt,
	)
	return err
}

func (m *Model) ConsumeAuthNonce(ctx context.Context, id string, now time.Time) (AuthNonce, error) {
	var nonce AuthNonce
	err := m.db.QueryRowContext(ctx, `
		UPDATE auth_nonces SET used_at=$2
		WHERE id=$1 AND used_at IS NULL AND expires_at > $2
		RETURNING id, chain_id, wallet, nonce_hash, message, expires_at`, id, now,
	).Scan(&nonce.ID, &nonce.ChainID, &nonce.Wallet, &nonce.NonceHash, &nonce.Message, &nonce.ExpiresAt)
	if err == sql.ErrNoRows {
		return AuthNonce{}, sql.ErrNoRows
	}
	return nonce, err
}
