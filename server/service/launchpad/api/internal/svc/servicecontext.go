package svc

import (
	"net/http"

	"o1-launchpad/service/launchpad/api/internal/config"
	"o1-launchpad/service/launchpad/api/internal/middleware"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Launchpad launchpadclient.Launchpad
	Session   func(http.HandlerFunc) http.HandlerFunc
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Launchpad: launchpadclient.NewLaunchpad(zrpc.MustNewClient(c.LaunchpadRpc)),
		Session:   middleware.SessionMiddleware(c.Auth.SessionSecret),
	}
}
