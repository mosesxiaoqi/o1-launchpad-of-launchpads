package logic

import (
	"context"

	"o1-launchpad/common/pagination"
	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	limit := pagination.Limit(in.Limit)
	cursor, err := pagination.DecodeOptional(in.Cursor)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	items, err := l.svcCtx.Model.ListLaunchpadTokensPage(l.ctx, in.ChainId, in.Slug, limit+1, cursor)
	if err != nil {
		return nil, err
	}
	response := &launchpad.ListTokensResponse{Tokens: make([]*launchpad.Token, 0, len(items))}
	if len(items) > limit {
		last := items[limit-1]
		response.NextCursor = pagination.Encode(pagination.Cursor{CreatedAt: last.CreatedAt, ID: last.ID})
		items = items[:limit]
	}
	for _, item := range items {
		response.Tokens = append(response.Tokens, tokenInfo(item))
	}
	return response, nil
}
