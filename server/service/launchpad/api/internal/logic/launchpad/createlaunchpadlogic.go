package launchpad

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/middleware"
	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLaunchpadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLaunchpadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLaunchpadLogic {
	return &CreateLaunchpadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLaunchpadLogic) CreateLaunchpad(req *types.CreateLaunchpadRequest) (resp *types.LaunchpadInfo, err error) {
	wallet, ok := middleware.AddressFromContext(l.ctx)
	if !ok {
		return nil, context.Canceled
	}
	result, err := l.svcCtx.Launchpad.CreateLaunchpad(l.ctx, &launchpadclient.CreateLaunchpadRequest{
		ChainId: req.ChainId, Slug: req.Slug, Name: req.Name, Description: req.Description,
		LogoUrl: req.LogoUrl, PrimaryColor: req.PrimaryColor, RegistryTxHash: req.RegistryTxHash, Wallet: wallet,
	})
	if err != nil {
		return nil, err
	}
	mapped := mapLaunchpad(result)
	return &mapped, nil
}
