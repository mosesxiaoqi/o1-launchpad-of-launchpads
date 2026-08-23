package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLaunchpadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLaunchpadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLaunchpadLogic {
	return &GetLaunchpadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLaunchpadLogic) GetLaunchpad(in *launchpad.GetLaunchpadRequest) (*launchpad.LaunchpadInfo, error) {
	// todo: add your logic here and delete this line

	return &launchpad.LaunchpadInfo{}, nil
}
