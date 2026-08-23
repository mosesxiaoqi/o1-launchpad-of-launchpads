package auth

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
