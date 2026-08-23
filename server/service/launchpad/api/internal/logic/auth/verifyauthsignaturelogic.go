package auth

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyAuthSignatureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyAuthSignatureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyAuthSignatureLogic {
	return &VerifyAuthSignatureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyAuthSignatureLogic) VerifyAuthSignature(req *types.VerifyAuthSignatureRequest) (resp *types.VerifyAuthSignatureResponse, err error) {
	result, err := l.svcCtx.Launchpad.VerifyAuthSignature(l.ctx, &launchpadclient.VerifyAuthSignatureRequest{
		ChallengeId: req.ChallengeId, Address: req.Address, Signature: req.Signature,
	})
	if err != nil {
		return nil, err
	}
	return &types.VerifyAuthSignatureResponse{Address: result.Address, ExpiresAt: result.ExpiresAt}, nil
}
