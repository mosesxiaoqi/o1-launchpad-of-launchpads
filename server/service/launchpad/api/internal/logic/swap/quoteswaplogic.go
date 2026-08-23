package swap

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
