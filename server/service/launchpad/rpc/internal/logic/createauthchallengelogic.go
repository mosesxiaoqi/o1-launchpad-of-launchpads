package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAuthChallengeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAuthChallengeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAuthChallengeLogic {
	return &CreateAuthChallengeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAuthChallengeLogic) CreateAuthChallenge(in *launchpad.CreateAuthChallengeRequest) (*launchpad.CreateAuthChallengeResponse, error) {
	// todo: add your logic here and delete this line

	return &launchpad.CreateAuthChallengeResponse{}, nil
}
