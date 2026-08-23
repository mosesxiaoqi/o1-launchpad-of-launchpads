package launchpad

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLaunchpadTokensLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLaunchpadTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLaunchpadTokensLogic {
	return &ListLaunchpadTokensLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLaunchpadTokensLogic) ListLaunchpadTokens(req *types.ListLaunchpadTokensRequest) (resp *types.ListTokensResponse, err error) {
	result, err := l.svcCtx.Launchpad.ListLaunchpadTokens(l.ctx, &launchpadclient.ListLaunchpadTokensRequest{
		ChainId: 84532, Slug: req.Slug, Limit: req.Limit, Cursor: req.Cursor,
	})
	if err != nil {
		return nil, err
	}
	return &types.ListTokensResponse{Tokens: mapTokens(result.Tokens), NextCursor: result.NextCursor}, nil
}
