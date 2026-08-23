package config

import (
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
)

func TestValidateRejectsMissingRPCOrSessionSecret(t *testing.T) {
	tests := []Config{
		{LaunchpadRpc: zrpc.RpcClientConf{Endpoints: nil}, Auth: Auth{SessionSecret: "secret"}},
		{LaunchpadRpc: zrpc.RpcClientConf{Endpoints: []string{"127.0.0.1:8080"}}, Auth: Auth{}},
	}
	for _, cfg := range tests {
		if err := cfg.Validate(); err == nil {
			t.Fatal("expected invalid gateway config")
		}
	}
}

func TestLoadExpandsSessionSecret(t *testing.T) {
	t.Setenv("SESSION_SECRET", "test-session-secret")
	var cfg Config
	if err := conf.Load("../../etc/launchpad-api.yaml", &cfg, conf.UseEnv()); err != nil {
		t.Fatal(err)
	}
	if cfg.Auth.SessionSecret != "test-session-secret" {
		t.Fatalf("SessionSecret = %q", cfg.Auth.SessionSecret)
	}
}

func TestValidateAcceptsGatewayConfig(t *testing.T) {
	cfg := Config{
		LaunchpadRpc: zrpc.RpcClientConf{Endpoints: []string{"127.0.0.1:8080"}},
		Auth:         Auth{SessionSecret: "secret", ChallengeTTLSeconds: 300},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
