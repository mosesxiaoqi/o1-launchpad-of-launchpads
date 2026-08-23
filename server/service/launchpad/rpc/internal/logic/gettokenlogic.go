package logic

import (
	"context"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTokenLogic) GetToken(in *launchpad.GetTokenRequest) (*launchpad.Token, error) {
	if !common.IsHexAddress(in.Token) {
		return nil, context.Canceled
	}
	item, err := l.svcCtx.Model.GetToken(l.ctx, in.ChainId, common.HexToAddress(in.Token).Bytes())
	if err != nil {
		return nil, err
	}
	return tokenInfo(item), nil
}
