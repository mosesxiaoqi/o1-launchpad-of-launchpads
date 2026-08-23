package system

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigLogic {
	return &GetConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConfigLogic) GetConfig() (resp *types.ConfigResponse, err error) {
	result, err := l.svcCtx.Launchpad.GetConfig(l.ctx, &launchpadclient.Empty{})
	if err != nil {
		return nil, err
	}
	return &types.ConfigResponse{
		ChainId: result.ChainId, Registry: result.Registry, Factory: result.Factory, Hook: result.Hook,
		FeeEscrow: result.FeeEscrow, Quote: result.Quote, ConfigVersion: result.ConfigVersion,
	}, nil
}
