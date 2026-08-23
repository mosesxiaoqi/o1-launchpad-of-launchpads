package config

import (
	"errors"
	"net"
	"net/url"
	"strings"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Auth struct {
	SessionSecret       string
	ChallengeTTLSeconds uint64
	AllowedOrigin       string
}

type Config struct {
	rest.RestConf
	LaunchpadRpc zrpc.RpcClientConf
	Auth         Auth
}

func (c Config) Validate() error {
	if len(c.LaunchpadRpc.Endpoints) == 0 {
		return errors.New("LaunchpadRpc.Endpoints is required")
	}
	if strings.TrimSpace(c.Auth.SessionSecret) == "" {
		return errors.New("Auth.SessionSecret is required")
	}
	if c.Auth.ChallengeTTLSeconds == 0 {
		return errors.New("Auth.ChallengeTTLSeconds must be positive")
	}
	if _, _, err := c.Auth.ChallengeOrigin(); err != nil {
		return err
	}
	return nil
}

func (a Auth) ChallengeOrigin() (domain string, origin string, err error) {
	parsed, err := url.Parse(a.AllowedOrigin)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" ||
		parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.User != nil {
		return "", "", errors.New("Auth.AllowedOrigin must be an HTTP(S) origin without a path")
	}
	if parsed.Scheme == "http" && !strings.EqualFold(parsed.Hostname(), "localhost") && !net.ParseIP(parsed.Hostname()).IsLoopback() {
		return "", "", errors.New("Auth.AllowedOrigin HTTP is allowed only for loopback development")
	}
	return parsed.Hostname(), parsed.String(), nil
}
