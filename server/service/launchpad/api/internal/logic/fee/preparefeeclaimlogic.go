package fee

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

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
	result, err := l.svcCtx.Launchpad.PrepareFeeClaim(l.ctx, &launchpadclient.PrepareFeeClaimRequest{
		ChainId: req.ChainId, Currency: req.Currency, Recipient: req.Recipient, Wallet: req.Wallet,
	})
	if err != nil {
		return nil, err
	}
	return &types.PreparedTransaction{
		ChainId: result.ChainId, From: result.From, To: result.To, Data: result.Data, Value: result.Value,
		Deadline: result.Deadline, Review: result.Review,
	}, nil
}
