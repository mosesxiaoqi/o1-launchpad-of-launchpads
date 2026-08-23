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
	Domain    string
	URI       string
	Nonce     string
	Message   string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func (m *Model) CreateAuthNonce(ctx context.Context, nonce AuthNonce) error {
	_, err := m.db.ExecContext(ctx, `
		INSERT INTO auth_nonces (id, chain_id, wallet, nonce_hash, domain, uri, nonce, message, issued_at, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		nonce.ID, nonce.ChainID, nonce.Wallet, nonce.NonceHash, nonce.Domain, nonce.URI,
		nonce.Nonce, nonce.Message, nonce.IssuedAt, nonce.ExpiresAt,
	)
	return err
}

func (m *Model) GetAuthNonce(ctx context.Context, id string) (AuthNonce, error) {
	var nonce AuthNonce
	err := m.db.QueryRowContext(ctx, `
		SELECT id, chain_id, wallet, nonce_hash, domain, uri, nonce, message, issued_at, expires_at
		FROM auth_nonces WHERE id=$1 AND used_at IS NULL`, id,
	).Scan(
		&nonce.ID, &nonce.ChainID, &nonce.Wallet, &nonce.NonceHash, &nonce.Domain, &nonce.URI,
		&nonce.Nonce, &nonce.Message, &nonce.IssuedAt, &nonce.ExpiresAt,
	)
	return nonce, err
}

func (m *Model) ConsumeAuthNonce(ctx context.Context, id string, now time.Time) (AuthNonce, error) {
	var nonce AuthNonce
	err := m.db.QueryRowContext(ctx, `
		UPDATE auth_nonces SET used_at=$2
		WHERE id=$1 AND used_at IS NULL AND expires_at > $2
		RETURNING id, chain_id, wallet, nonce_hash, domain, uri, nonce, message, issued_at, expires_at`, id, now,
	).Scan(
		&nonce.ID, &nonce.ChainID, &nonce.Wallet, &nonce.NonceHash, &nonce.Domain, &nonce.URI,
		&nonce.Nonce, &nonce.Message, &nonce.IssuedAt, &nonce.ExpiresAt,
	)
	if err == sql.ErrNoRows {
		return AuthNonce{}, sql.ErrNoRows
	}
	return nonce, err
}
