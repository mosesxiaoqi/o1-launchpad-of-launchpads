package launch

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PrepareLaunchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPrepareLaunchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareLaunchLogic {
	return &PrepareLaunchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PrepareLaunchLogic) PrepareLaunch(req *types.PrepareLaunchRequest) (resp *types.PreparedTransaction, err error) {
	// todo: add your logic here and delete this line

	return
}
