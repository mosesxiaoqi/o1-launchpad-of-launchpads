package swap

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuoteSwapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuoteSwapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuoteSwapLogic {
	return &QuoteSwapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuoteSwapLogic) QuoteSwap(req *types.QuoteSwapRequest) (resp *types.QuoteSwapResponse, err error) {
	result, err := l.svcCtx.Launchpad.QuoteSwap(l.ctx, &launchpadclient.QuoteSwapRequest{
		ChainId: req.ChainId, Token: req.Token, Amount: req.Amount, Buy: req.Buy, Wallet: req.Wallet, Referrer: req.Referrer,
	})
	if err != nil {
		return nil, err
	}
	return &types.QuoteSwapResponse{
		QuoteId: result.QuoteId, AmountIn: result.AmountIn, AmountOut: result.AmountOut,
		Fee: result.Fee, ExpiresAt: result.ExpiresAt,
	}, nil
}
