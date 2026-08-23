package token

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenLogic) GetToken(req *types.TokenPath) (resp *types.TokenInfo, err error) {
	item, err := l.svcCtx.Launchpad.GetToken(l.ctx, &launchpadclient.GetTokenRequest{ChainId: req.ChainId, Token: req.TokenAddress})
	if err != nil {
		return nil, err
	}
	return &types.TokenInfo{
		ChainId: item.ChainId, Token: item.Token, PoolId: item.PoolId, LaunchpadId: item.LaunchpadId,
		LaunchpadSlug: item.LaunchpadSlug, Creator: item.Creator, Quote: item.Quote, Supply: item.Supply,
		TxHash: item.TxHash, BlockNumber: item.BlockNumber, Status: item.Status, CreatedAt: item.CreatedAt,
	}, nil
}
