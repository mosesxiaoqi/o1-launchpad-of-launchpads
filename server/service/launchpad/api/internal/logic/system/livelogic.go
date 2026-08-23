package system

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LiveLogic {
	return &LiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LiveLogic) Live() (resp *types.HealthResponse, err error) {
	return &types.HealthResponse{Status: "ok"}, nil
}
