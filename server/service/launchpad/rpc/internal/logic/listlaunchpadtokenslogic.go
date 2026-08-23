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
	// todo: add your logic here and delete this line

	return &launchpad.ListTokensResponse{}, nil
}
