package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type PrepareFeeClaimLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPrepareFeeClaimLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareFeeClaimLogic {
	return &PrepareFeeClaimLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PrepareFeeClaimLogic) PrepareFeeClaim(in *launchpad.PrepareFeeClaimRequest) (*launchpad.PreparedTransaction, error) {
	// todo: add your logic here and delete this line

	return &launchpad.PreparedTransaction{}, nil
}
