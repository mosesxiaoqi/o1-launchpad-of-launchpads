package config

import (
	"testing"

	commonconfig "o1-launchpad/common/config"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestValidateRejectsMissingDatabaseAndWrongChain(t *testing.T) {
	tests := []Config{
		{Chain: commonconfig.Chain{ChainId: 84532}, DeploymentFile: "deployment.json"},
		{
			Database:       commonconfig.Database{DataSource: "postgres://db"},
			Chain:          commonconfig.Chain{ChainId: 1, HttpRpc: "http://rpc", Confirmations: 2},
			DeploymentFile: "deployment.json",
		},
	}
	for _, cfg := range tests {
		if err := cfg.Validate(); err == nil {
			t.Fatal("expected invalid rpc config")
		}
	}
}

func TestLoadExpandsDatabaseAndRPC(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test-db")
	t.Setenv("BASE_SEPOLIA_RPC_URL", "https://base-sepolia.example")
	t.Setenv("BASE_SEPOLIA_WS_URL", "wss://base-sepolia.example")
	var cfg Config
	if err := conf.Load("../../etc/launchpad-rpc.yaml", &cfg, conf.UseEnv()); err != nil {
		t.Fatal(err)
	}
	if cfg.Database.DataSource != "postgres://test-db" || cfg.Chain.HttpRpc != "https://base-sepolia.example" {
		t.Fatalf("environment was not expanded: %#v", cfg)
	}
}

func TestValidateAcceptsRPCConfig(t *testing.T) {
	cfg := Config{
		Database:       commonconfig.Database{DataSource: "postgres://db"},
		Chain:          commonconfig.Chain{ChainId: 84532, HttpRpc: "http://rpc", Confirmations: 2},
		DeploymentFile: "deployment.json",
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
