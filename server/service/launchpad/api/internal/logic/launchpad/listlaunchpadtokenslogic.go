package launchpad

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
