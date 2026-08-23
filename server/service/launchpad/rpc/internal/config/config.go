package config

import (
	"errors"
	"strings"

	commonconfig "o1-launchpad/common/config"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Database       commonconfig.Database
	Chain          commonconfig.Chain
	DeploymentFile string
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Database.DataSource) == "" {
		return errors.New("Database.DataSource is required")
	}
	if c.Chain.ChainId != 84532 {
		return errors.New("Chain.ChainId must be Base Sepolia (84532)")
	}
	if strings.TrimSpace(c.Chain.HttpRpc) == "" {
		return errors.New("Chain.HttpRpc is required")
	}
	if c.Chain.Confirmations == 0 {
		return errors.New("Chain.Confirmations must be positive")
	}
	if strings.TrimSpace(c.DeploymentFile) == "" {
		return errors.New("DeploymentFile is required")
	}
	return nil
}
