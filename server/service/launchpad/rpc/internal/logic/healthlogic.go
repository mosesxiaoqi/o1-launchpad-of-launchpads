package logic

import (
	"context"
	"errors"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthLogic {
	return &HealthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HealthLogic) Health(in *launchpad.Empty) (*launchpad.HealthResponse, error) {
	if err := l.svcCtx.DB.PingContext(l.ctx); err != nil {
		return nil, err
	}
	chainID, err := l.svcCtx.Chain.RPC.ChainID(l.ctx)
	if err != nil {
		return nil, err
	}
	if chainID.Int64() != l.svcCtx.Config.Chain.ChainId {
		return nil, errors.New("connected chain id does not match configuration")
	}
	return &launchpad.HealthResponse{Status: "ok"}, nil
}
