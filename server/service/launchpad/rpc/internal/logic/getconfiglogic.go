package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigLogic {
	return &GetConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConfigLogic) GetConfig(in *launchpad.Empty) (*launchpad.GetConfigResponse, error) {
	deployment := l.svcCtx.Deployment
	return &launchpad.GetConfigResponse{
		ChainId: deployment.ChainID, Registry: deployment.Registry, Factory: deployment.Factory,
		Hook: deployment.Hook, FeeEscrow: deployment.FeeEscrow, Quote: deployment.Quote,
		ConfigVersion: 1, LaasTreasury: deployment.LaaSTreasury,
	}, nil
}
