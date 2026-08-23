package swap

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

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
	result, err := l.svcCtx.Launchpad.PrepareSwap(l.ctx, &launchpadclient.PrepareSwapRequest{
		ChainId: req.ChainId, QuoteId: req.QuoteId, SlippageBps: req.SlippageBps, Wallet: req.Wallet,
	})
	if err != nil {
		return nil, err
	}
	return &types.PreparedTransaction{
		ChainId: result.ChainId, From: result.From, To: result.To, Data: result.Data, Value: result.Value,
		Deadline: result.Deadline, Review: result.Review,
	}, nil
}
