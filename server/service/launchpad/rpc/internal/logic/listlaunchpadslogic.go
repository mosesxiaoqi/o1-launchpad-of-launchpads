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
	limit := int(in.Limit)
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	items, err := l.svcCtx.Model.ListLaunchpads(l.ctx, limit)
	if err != nil {
		return nil, err
	}
	response := &launchpad.ListLaunchpadsResponse{Launchpads: make([]*launchpad.LaunchpadInfo, 0, len(items))}
	for _, item := range items {
		response.Launchpads = append(response.Launchpads, launchpadInfo(item))
	}
	return response, nil
}
