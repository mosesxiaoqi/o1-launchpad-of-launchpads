package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyAuthSignatureLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyAuthSignatureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyAuthSignatureLogic {
	return &VerifyAuthSignatureLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyAuthSignatureLogic) VerifyAuthSignature(in *launchpad.VerifyAuthSignatureRequest) (*launchpad.VerifyAuthSignatureResponse, error) {
	// todo: add your logic here and delete this line

	return &launchpad.VerifyAuthSignatureResponse{}, nil
}
