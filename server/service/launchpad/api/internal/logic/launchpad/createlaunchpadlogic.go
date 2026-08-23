package launchpad

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
