package model

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.PingContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	_, filename, _, _ := runtime.Caller(0)
	migration, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "..", "..", "migrations", "001_init.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(migration)); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func cleanTables(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`TRUNCATE auth_nonces, token_launches, chain_transactions, indexer_checkpoints, launchpads RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestApplyLaunchEventIsIdempotentAndRejectsOwnershipConflict(t *testing.T) {
	db := openTestDB(t)
	cleanTables(t, db)
	m := New(db)
	ctx := context.Background()
	launchpad := sampleLaunchpad()
	if err := m.CreateLaunchpad(ctx, launchpad); err != nil {
		t.Fatal(err)
	}

	event := sampleLaunchEvent()
	checkpoint := Checkpoint{ChainID: event.ChainID, Factory: event.Factory, NextBlock: event.BlockNumber + 1, LastBlockHash: event.BlockHash}
	if err := m.ApplyLaunchEvent(ctx, event, checkpoint); err != nil {
		t.Fatal(err)
	}
	if err := m.ApplyLaunchEvent(ctx, event, checkpoint); err != nil {
		t.Fatalf("idempotent replay failed: %v", err)
	}

	var count int
	if err := db.QueryRow(`SELECT count(*) FROM token_launches`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("token launch count = %d, want 1", count)
	}

	conflict := event
	conflict.PoolID = bytesOf(0x99, 32)
	if err := m.ApplyLaunchEvent(ctx, conflict, checkpoint); !errors.Is(err, ErrOwnershipConflict) {
		t.Fatalf("conflict error = %v, want ErrOwnershipConflict", err)
	}
}

func TestTransactionLifecycleAndOrphanRemoval(t *testing.T) {
	db := openTestDB(t)
	cleanTables(t, db)
	m := New(db)
	ctx := context.Background()
	if err := m.CreateLaunchpad(ctx, sampleLaunchpad()); err != nil {
		t.Fatal(err)
	}
	event := sampleLaunchEvent()

	for _, status := range []string{"pending", "confirming", "confirmed"} {
		if err := m.UpsertTransaction(ctx, ChainTransaction{
			ChainID: event.ChainID, TxHash: event.TxHash, Kind: "launch", Status: status,
		}); err != nil {
			t.Fatal(err)
		}
	}
	tx, err := m.GetTransaction(ctx, event.ChainID, event.TxHash)
	if err != nil {
		t.Fatal(err)
	}
	if tx.Status != "confirmed" {
		t.Fatalf("transaction status = %q", tx.Status)
	}

	cp := Checkpoint{ChainID: event.ChainID, Factory: event.Factory, NextBlock: event.BlockNumber + 1, LastBlockHash: event.BlockHash}
	if err := m.ApplyLaunchEvent(ctx, event, cp); err != nil {
		t.Fatal(err)
	}
	if err := m.DeleteLaunchEventsFromBlock(ctx, event.ChainID, event.Factory, event.BlockNumber); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := db.QueryRow(`SELECT count(*) FROM token_launches`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("orphan launch count = %d", count)
	}
	got, err := m.GetCheckpoint(ctx, event.ChainID, event.Factory)
	if err != nil {
		t.Fatal(err)
	}
	if got.NextBlock != event.BlockNumber || len(got.LastBlockHash) != 0 {
		t.Fatalf("checkpoint after reorg = %#v", got)
	}
}

func sampleLaunchpad() Launchpad {
	return Launchpad{
		ChainID: 84532, LaunchpadID: bytesOf(0x11, 32), Slug: "ai-pad", Name: "AI Pad",
		Owner: bytesOf(0x22, 20), Treasury: bytesOf(0x33, 20), RegistryTxHash: bytesOf(0x44, 32),
		RegistryBlockNumber: 100, Active: true, CreatedAt: time.Unix(1_700_000_000, 0).UTC(),
	}
}

func sampleLaunchEvent() LaunchEvent {
	return LaunchEvent{
		ChainID: 84532, LaunchpadID: bytesOf(0x11, 32), Token: bytesOf(0x55, 20),
		PoolID: bytesOf(0x66, 32), Creator: bytesOf(0x77, 20), Quote: bytesOf(0x88, 20),
		Factory: bytesOf(0xaa, 20), Supply: "1000000000000000000000000000", TxHash: bytesOf(0xbb, 32),
		BlockNumber: 120, BlockHash: bytesOf(0xcc, 32), LogIndex: 3, Status: "confirmed",
		CreatedAt: time.Unix(1_700_000_100, 0).UTC(),
	}
}

func bytesOf(value byte, size int) []byte {
	out := make([]byte, size)
	for i := range out {
		out[i] = value
	}
	return out
}
