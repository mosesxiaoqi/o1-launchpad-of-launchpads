package launch

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PrepareLaunchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPrepareLaunchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareLaunchLogic {
	return &PrepareLaunchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PrepareLaunchLogic) PrepareLaunch(req *types.PrepareLaunchRequest) (resp *types.PreparedTransaction, err error) {
	result, err := l.svcCtx.Launchpad.PrepareLaunch(l.ctx, &launchpadclient.PrepareLaunchRequest{
		ChainId: req.ChainId, LaunchpadId: req.LaunchpadId, Name: req.Name, Symbol: req.Symbol,
		ContractUri: req.ContractUri, Salt: req.Salt, Wallet: req.Wallet,
	})
	if err != nil {
		return nil, err
	}
	return &types.PreparedTransaction{
		ChainId: result.ChainId, From: result.From, To: result.To, Data: result.Data, Value: result.Value,
		Deadline: result.Deadline, Review: result.Review,
	}, nil
}
