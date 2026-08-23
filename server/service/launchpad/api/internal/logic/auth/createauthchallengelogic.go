package auth

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
