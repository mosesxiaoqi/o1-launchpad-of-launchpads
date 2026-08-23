package launchpad

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLaunchpadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLaunchpadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLaunchpadLogic {
	return &GetLaunchpadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLaunchpadLogic) GetLaunchpad(req *types.LaunchpadPath) (resp *types.LaunchpadInfo, err error) {
	item, err := l.svcCtx.Launchpad.GetLaunchpad(l.ctx, &launchpadclient.GetLaunchpadRequest{ChainId: 84532, Slug: req.Slug})
	if err != nil {
		return nil, err
	}
	mapped := mapLaunchpad(item)
	return &mapped, nil
}
