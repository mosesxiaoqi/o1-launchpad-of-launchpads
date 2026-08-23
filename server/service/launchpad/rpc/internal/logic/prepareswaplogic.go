package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type PrepareSwapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPrepareSwapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareSwapLogic {
	return &PrepareSwapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PrepareSwapLogic) PrepareSwap(in *launchpad.PrepareSwapRequest) (*launchpad.PreparedTransaction, error) {
	// todo: add your logic here and delete this line

	return &launchpad.PreparedTransaction{}, nil
}
