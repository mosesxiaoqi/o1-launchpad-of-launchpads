package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTokensLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTokensLogic {
	return &ListTokensLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListTokensLogic) ListTokens(in *launchpad.ListTokensRequest) (*launchpad.ListTokensResponse, error) {
	limit := int(in.Limit)
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	items, err := l.svcCtx.Model.ListTokens(l.ctx, in.ChainId, limit)
	if err != nil {
		return nil, err
	}
	response := &launchpad.ListTokensResponse{Tokens: make([]*launchpad.Token, 0, len(items))}
	for _, item := range items {
		response.Tokens = append(response.Tokens, tokenInfo(item))
	}
	return response, nil
}
