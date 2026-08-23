package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLaunchpadsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLaunchpadsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLaunchpadsLogic {
	return &ListLaunchpadsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLaunchpadsLogic) ListLaunchpads(in *launchpad.PageRequest) (*launchpad.ListLaunchpadsResponse, error) {
	// todo: add your logic here and delete this line

	return &launchpad.ListLaunchpadsResponse{}, nil
}
