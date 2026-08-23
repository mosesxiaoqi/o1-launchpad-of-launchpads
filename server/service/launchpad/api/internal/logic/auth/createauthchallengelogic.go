package auth

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAuthChallengeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAuthChallengeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAuthChallengeLogic {
	return &CreateAuthChallengeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAuthChallengeLogic) CreateAuthChallenge(req *types.CreateAuthChallengeRequest) (resp *types.CreateAuthChallengeResponse, err error) {
	domain, origin, err := l.svcCtx.Config.Auth.ChallengeOrigin()
	if err != nil {
		return nil, err
	}
	result, err := l.svcCtx.Launchpad.CreateAuthChallenge(l.ctx, &launchpadclient.CreateAuthChallengeRequest{
		ChainId: req.ChainId, Address: req.Address, Domain: domain, Uri: origin,
	})
	if err != nil {
		return nil, err
	}
	return &types.CreateAuthChallengeResponse{
		ChallengeId: result.ChallengeId, Message: result.Message, ExpiresAt: result.ExpiresAt,
	}, nil
}
