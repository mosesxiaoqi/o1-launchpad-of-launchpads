package config

import (
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
)

func TestValidateRejectsMissingRPCOrSessionSecret(t *testing.T) {
	tests := []Config{
		{LaunchpadRpc: zrpc.RpcClientConf{Endpoints: nil}, Auth: Auth{SessionSecret: "secret", AllowedOrigin: "http://localhost:3000"}},
		{LaunchpadRpc: zrpc.RpcClientConf{Endpoints: []string{"127.0.0.1:8080"}}, Auth: Auth{}},
		{LaunchpadRpc: zrpc.RpcClientConf{Endpoints: []string{"127.0.0.1:8080"}}, Auth: Auth{SessionSecret: "secret", ChallengeTTLSeconds: 300}},
		{LaunchpadRpc: zrpc.RpcClientConf{Endpoints: []string{"127.0.0.1:8080"}}, Auth: Auth{SessionSecret: "secret", ChallengeTTLSeconds: 300, AllowedOrigin: "http://demo.example"}},
	}
	for _, cfg := range tests {
		if err := cfg.Validate(); err == nil {
			t.Fatal("expected invalid gateway config")
		}
	}
}

func TestLoadExpandsSessionSecret(t *testing.T) {
	t.Setenv("SESSION_SECRET", "test-session-secret")
	t.Setenv("FRONTEND_ORIGIN", "http://localhost:3000")
	t.Setenv("SERVICE_MODE", "pro")
	var cfg Config
	if err := conf.Load("../../etc/launchpad-api.yaml", &cfg, conf.UseEnv()); err != nil {
		t.Fatal(err)
	}
	if cfg.Auth.SessionSecret != "test-session-secret" {
		t.Fatalf("SessionSecret = %q", cfg.Auth.SessionSecret)
	}
	if cfg.Auth.AllowedOrigin != "http://localhost:3000" {
		t.Fatalf("AllowedOrigin = %q", cfg.Auth.AllowedOrigin)
	}
	if cfg.Host != "127.0.0.1" || cfg.Mode != "pro" {
		t.Fatalf("gateway exposure config = host %q, mode %q", cfg.Host, cfg.Mode)
	}
}

func TestValidateAcceptsGatewayConfig(t *testing.T) {
	cfg := Config{
		LaunchpadRpc: zrpc.RpcClientConf{Endpoints: []string{"127.0.0.1:8080"}},
		Auth:         Auth{SessionSecret: "secret", ChallengeTTLSeconds: 300, AllowedOrigin: "http://localhost:3000"},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestAuthChallengeOriginUsesTrustedConfiguration(t *testing.T) {
	auth := Auth{AllowedOrigin: "https://demo.example:8443"}
	domain, origin, err := auth.ChallengeOrigin()
	if err != nil {
		t.Fatal(err)
	}
	if domain != "demo.example" || origin != "https://demo.example:8443" {
		t.Fatalf("ChallengeOrigin() = %q, %q", domain, origin)
	}
}
