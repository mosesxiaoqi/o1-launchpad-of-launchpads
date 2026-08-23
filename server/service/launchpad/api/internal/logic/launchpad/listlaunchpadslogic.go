package launchpad

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLaunchpadsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLaunchpadsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLaunchpadsLogic {
	return &ListLaunchpadsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLaunchpadsLogic) ListLaunchpads(req *types.PageRequest) (resp *types.ListLaunchpadsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
