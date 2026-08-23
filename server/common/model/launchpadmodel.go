package model

import (
	"context"
	"database/sql"
	"time"
)

type Launchpad struct {
	ID                  int64
	ChainID             int64
	LaunchpadID         []byte
	Slug                string
	Name                string
	Description         string
	LogoURL             string
	PrimaryColor        string
	Owner               []byte
	Treasury            []byte
	RegistryTxHash      []byte
	RegistryBlockNumber uint64
	Active              bool
	CreatedAt           time.Time
}

func (m *Model) CreateLaunchpad(ctx context.Context, launchpad Launchpad) error {
	_, err := m.db.ExecContext(ctx, `
		INSERT INTO launchpads (
			chain_id, launchpad_id, slug, name, description, logo_url, primary_color,
			owner, treasury, registry_tx_hash, registry_block_number, active, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		launchpad.ChainID, launchpad.LaunchpadID, launchpad.Slug, launchpad.Name, launchpad.Description,
		launchpad.LogoURL, launchpad.PrimaryColor, launchpad.Owner, launchpad.Treasury,
		launchpad.RegistryTxHash, launchpad.RegistryBlockNumber, launchpad.Active, launchpad.CreatedAt,
	)
	return err
}

func (m *Model) GetLaunchpadBySlug(ctx context.Context, chainID int64, slug string) (Launchpad, error) {
	var launchpad Launchpad
	err := m.db.QueryRowContext(ctx, `
		SELECT id, chain_id, launchpad_id, slug, name, description, logo_url, primary_color,
		       owner, treasury, registry_tx_hash, registry_block_number, active, created_at
		FROM launchpads WHERE chain_id=$1 AND slug=$2`, chainID, slug,
	).Scan(
		&launchpad.ID, &launchpad.ChainID, &launchpad.LaunchpadID, &launchpad.Slug, &launchpad.Name,
		&launchpad.Description, &launchpad.LogoURL, &launchpad.PrimaryColor, &launchpad.Owner,
		&launchpad.Treasury, &launchpad.RegistryTxHash, &launchpad.RegistryBlockNumber,
		&launchpad.Active, &launchpad.CreatedAt,
	)
	return launchpad, err
}

func (m *Model) ListLaunchpads(ctx context.Context, limit int) ([]Launchpad, error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT id, chain_id, launchpad_id, slug, name, description, logo_url, primary_color,
		       owner, treasury, registry_tx_hash, registry_block_number, active, created_at
		FROM launchpads ORDER BY created_at DESC, id DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var launchpads []Launchpad
	for rows.Next() {
		var launchpad Launchpad
		if err := rows.Scan(
			&launchpad.ID, &launchpad.ChainID, &launchpad.LaunchpadID, &launchpad.Slug, &launchpad.Name,
			&launchpad.Description, &launchpad.LogoURL, &launchpad.PrimaryColor, &launchpad.Owner,
			&launchpad.Treasury, &launchpad.RegistryTxHash, &launchpad.RegistryBlockNumber,
			&launchpad.Active, &launchpad.CreatedAt,
		); err != nil {
			return nil, err
		}
		launchpads = append(launchpads, launchpad)
	}
	return launchpads, rows.Err()
}

func launchpadPrimaryKey(ctx context.Context, tx *sql.Tx, chainID int64, launchpadID []byte) (int64, error) {
	var id int64
	err := tx.QueryRowContext(ctx,
		`SELECT id FROM launchpads WHERE chain_id=$1 AND launchpad_id=$2`, chainID, launchpadID,
	).Scan(&id)
	return id, err
}
