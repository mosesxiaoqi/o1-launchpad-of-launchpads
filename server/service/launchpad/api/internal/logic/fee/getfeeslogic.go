package fee

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFeesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeesLogic {
	return &GetFeesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFeesLogic) GetFees(req *types.GetFeesRequest) (resp *types.GetFeesResponse, err error) {
	result, err := l.svcCtx.Launchpad.GetFees(l.ctx, &launchpadclient.GetFeesRequest{
		ChainId: 84532, Currency: req.Currency, Recipient: req.Recipient,
	})
	if err != nil {
		return nil, err
	}
	return &types.GetFeesResponse{
		ChainId: result.ChainId, Currency: result.Currency, Recipient: result.Recipient,
		Amount: result.Amount, BlockNumber: result.BlockNumber,
	}, nil
}
