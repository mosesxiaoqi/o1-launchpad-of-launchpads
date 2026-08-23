package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLaunchpadTokensLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLaunchpadTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLaunchpadTokensLogic {
	return &ListLaunchpadTokensLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLaunchpadTokensLogic) ListLaunchpadTokens(in *launchpad.ListLaunchpadTokensRequest) (*launchpad.ListTokensResponse, error) {
	limit := int(in.Limit)
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	items, err := l.svcCtx.Model.ListLaunchpadTokens(l.ctx, in.ChainId, in.Slug, limit)
	if err != nil {
		return nil, err
	}
	response := &launchpad.ListTokensResponse{Tokens: make([]*launchpad.Token, 0, len(items))}
	for _, item := range items {
		response.Tokens = append(response.Tokens, tokenInfo(item))
	}
	return response, nil
}
