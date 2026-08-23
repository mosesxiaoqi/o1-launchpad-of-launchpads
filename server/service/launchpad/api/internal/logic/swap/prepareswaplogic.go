package swap

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PrepareSwapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPrepareSwapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareSwapLogic {
	return &PrepareSwapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PrepareSwapLogic) PrepareSwap(req *types.PrepareSwapRequest) (resp *types.PreparedTransaction, err error) {
	// todo: add your logic here and delete this line

	return
}
