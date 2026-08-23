package config

import (
	"errors"
	"strings"

	commonconfig "o1-launchpad/common/config"

	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	service.ServiceConf
	Database       commonconfig.Database
	Chain          commonconfig.Chain
	DeploymentFile string
	Indexer        commonconfig.Indexer
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Database.DataSource) == "" {
		return errors.New("Database.DataSource is required")
	}
	if c.Chain.ChainId != 84532 || strings.TrimSpace(c.Chain.HttpRpc) == "" {
		return errors.New("Chain must configure Base Sepolia HTTP RPC")
	}
	if c.Chain.Confirmations == 0 {
		return errors.New("Chain.Confirmations must be positive")
	}
	if c.Indexer.BatchSize == 0 || c.Indexer.BatchSize > 500 {
		return errors.New("Indexer.BatchSize must be between 1 and 500")
	}
	if c.Indexer.PollInterval <= 0 || c.Indexer.ReorgLookback == 0 {
		return errors.New("Indexer poll interval and reorg lookback must be positive")
	}
	if strings.TrimSpace(c.DeploymentFile) == "" {
		return errors.New("DeploymentFile is required")
	}
	return nil
}
