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
	// todo: add your logic here and delete this line

	return &launchpad.ListTokensResponse{}, nil
}
