package token

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTokensLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTokensLogic {
	return &ListTokensLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTokensLogic) ListTokens(req *types.PageRequest) (resp *types.ListTokensResponse, err error) {
	result, err := l.svcCtx.Launchpad.ListTokens(l.ctx, &launchpadclient.ListTokensRequest{
		ChainId: 84532, Limit: req.Limit, Cursor: req.Cursor,
	})
	if err != nil {
		return nil, err
	}
	items := make([]types.TokenInfo, 0, len(result.Tokens))
	for _, item := range result.Tokens {
		items = append(items, types.TokenInfo{
			ChainId: item.ChainId, Token: item.Token, PoolId: item.PoolId, LaunchpadId: item.LaunchpadId,
			LaunchpadSlug: item.LaunchpadSlug, Creator: item.Creator, Quote: item.Quote, Supply: item.Supply,
			TxHash: item.TxHash, BlockNumber: item.BlockNumber, Status: item.Status, CreatedAt: item.CreatedAt,
		})
	}
	return &types.ListTokensResponse{Tokens: items, NextCursor: result.NextCursor}, nil
}
