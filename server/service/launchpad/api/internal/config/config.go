package config

import (
	"errors"
	"strings"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Auth struct {
	SessionSecret       string
	ChallengeTTLSeconds uint64
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
	return nil
}
