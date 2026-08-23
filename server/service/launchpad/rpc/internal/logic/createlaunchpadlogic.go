package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLaunchpadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLaunchpadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLaunchpadLogic {
	return &CreateLaunchpadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLaunchpadLogic) CreateLaunchpad(in *launchpad.CreateLaunchpadRequest) (*launchpad.LaunchpadInfo, error) {
	// todo: add your logic here and delete this line

	return &launchpad.LaunchpadInfo{}, nil
}
