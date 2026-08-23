package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type PrepareLaunchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPrepareLaunchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareLaunchLogic {
	return &PrepareLaunchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PrepareLaunchLogic) PrepareLaunch(in *launchpad.PrepareLaunchRequest) (*launchpad.PreparedTransaction, error) {
	// todo: add your logic here and delete this line

	return &launchpad.PreparedTransaction{}, nil
}
