package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuoteSwapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuoteSwapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuoteSwapLogic {
	return &QuoteSwapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QuoteSwapLogic) QuoteSwap(in *launchpad.QuoteSwapRequest) (*launchpad.QuoteSwapResponse, error) {
	// todo: add your logic here and delete this line

	return &launchpad.QuoteSwapResponse{}, nil
}
