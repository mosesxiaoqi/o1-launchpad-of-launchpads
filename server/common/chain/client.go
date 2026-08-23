package chain

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	commonconfig "o1-launchpad/common/config"
)

type Deployment struct {
	ChainID         int64  `json:"chainId"`
	PoolManager     string `json:"poolManager"`
	UniversalRouter string `json:"universalRouter"`
	StateView       string `json:"stateView"`
	Quoter          string `json:"quoter"`
	Permit2         string `json:"permit2"`
	Registry        string `json:"registry"`
	Factory         string `json:"factory"`
	Hook            string `json:"hook"`
	FeeEscrow       string `json:"feeEscrow"`
	Quote           string `json:"quote"`
	DeploymentBlock uint64 `json:"deploymentBlock"`
}

type Client struct {
	Chain      commonconfig.Chain
	Deployment Deployment
	Timeout    time.Duration
}

func LoadDeployment(path string, chainID int64) (Deployment, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Deployment{}, err
	}
	var deployment Deployment
	if err := json.Unmarshal(data, &deployment); err != nil {
		return Deployment{}, err
	}
	if deployment.ChainID != chainID || strings.TrimSpace(deployment.PoolManager) == "" {
		return Deployment{}, fmt.Errorf("deployment does not match chain %d", chainID)
	}
	return deployment, nil
}

func NewClient(config commonconfig.Chain, deployment Deployment) *Client {
	return &Client{Chain: config, Deployment: deployment, Timeout: 15 * time.Second}
}
