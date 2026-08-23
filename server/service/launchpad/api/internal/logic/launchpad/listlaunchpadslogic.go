package launchpad

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

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
	result, err := l.svcCtx.Launchpad.ListLaunchpads(l.ctx, &launchpadclient.PageRequest{Limit: req.Limit, Cursor: req.Cursor})
	if err != nil {
		return nil, err
	}
	response := &types.ListLaunchpadsResponse{NextCursor: result.NextCursor, Launchpads: make([]types.LaunchpadInfo, 0, len(result.Launchpads))}
	for _, item := range result.Launchpads {
		response.Launchpads = append(response.Launchpads, mapLaunchpad(item))
	}
	return response, nil
}
