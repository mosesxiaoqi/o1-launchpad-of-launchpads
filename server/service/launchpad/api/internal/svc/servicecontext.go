package svc

import (
	"o1-launchpad/service/launchpad/api/internal/config"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Launchpad launchpadclient.Launchpad
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Launchpad: launchpadclient.NewLaunchpad(zrpc.MustNewClient(c.LaunchpadRpc)),
	}
}
