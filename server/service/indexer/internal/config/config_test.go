package config

import (
	"testing"
	"time"

	commonconfig "o1-launchpad/common/config"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestValidateRejectsUnsafeScanConfig(t *testing.T) {
	base := Config{
		Database:       commonconfig.Database{DataSource: "postgres://db"},
		Chain:          commonconfig.Chain{ChainId: 84532, HttpRpc: "http://rpc", Confirmations: 2},
		DeploymentFile: "deployment.json",
		Indexer:        commonconfig.Indexer{BatchSize: 500, PollInterval: time.Second, ReorgLookback: 64},
	}

	badConfirmations := base
	badConfirmations.Chain.Confirmations = 0
	if err := badConfirmations.Validate(); err == nil {
		t.Fatal("expected confirmations validation error")
	}

	badBatch := base
	badBatch.Indexer.BatchSize = 501
	if err := badBatch.Validate(); err == nil {
		t.Fatal("expected batch size validation error")
	}
}

func TestLoadExpandsIndexerEnvironment(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://indexer-db")
	t.Setenv("BASE_SEPOLIA_RPC_URL", "https://base-sepolia.example")
	t.Setenv("BASE_SEPOLIA_WS_URL", "wss://base-sepolia.example")
	var cfg Config
	if err := conf.Load("../../etc/indexer.yaml", &cfg, conf.UseEnv()); err != nil {
		t.Fatal(err)
	}
	if cfg.Database.DataSource != "postgres://indexer-db" || cfg.Indexer.BatchSize != 10 {
		t.Fatalf("unexpected loaded config: %#v", cfg)
	}
}

func TestValidateAcceptsIndexerConfig(t *testing.T) {
	cfg := Config{
		Database:       commonconfig.Database{DataSource: "postgres://db"},
		Chain:          commonconfig.Chain{ChainId: 84532, HttpRpc: "http://rpc", Confirmations: 2},
		DeploymentFile: "deployment.json",
		Indexer:        commonconfig.Indexer{BatchSize: 500, PollInterval: 2 * time.Second, ReorgLookback: 64},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
