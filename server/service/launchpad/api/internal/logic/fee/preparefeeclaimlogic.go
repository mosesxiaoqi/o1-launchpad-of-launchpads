package fee

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PrepareFeeClaimLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPrepareFeeClaimLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareFeeClaimLogic {
	return &PrepareFeeClaimLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PrepareFeeClaimLogic) PrepareFeeClaim(req *types.PrepareFeeClaimRequest) (resp *types.PreparedTransaction, err error) {
	// todo: add your logic here and delete this line

	return
}
