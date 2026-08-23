package logic

import (
	"context"
	"errors"

	"o1-launchpad/common/chain"
	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFeesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeesLogic {
	return &GetFeesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFeesLogic) GetFees(in *launchpad.GetFeesRequest) (*launchpad.GetFeesResponse, error) {
	if in.ChainId != l.svcCtx.Config.Chain.ChainId || !common.IsHexAddress(in.Recipient) ||
		!common.IsHexAddress(in.Currency) || !common.IsHexAddress(l.svcCtx.Deployment.FeeEscrow) {
		return nil, errors.New("invalid fee query")
	}
	recipient := common.HexToAddress(in.Recipient)
	currency := common.HexToAddress(in.Currency)
	amount, blockNumber, err := chain.ReadOwedAtLatestBlock(
		l.ctx, l.svcCtx.Chain.RPC, common.HexToAddress(l.svcCtx.Deployment.FeeEscrow), recipient, currency,
	)
	if err != nil {
		return nil, err
	}
	return &launchpad.GetFeesResponse{
		ChainId: in.ChainId, Currency: currency.Hex(), Recipient: recipient.Hex(), Amount: amount.String(),
		BlockNumber: blockNumber,
	}, nil
}
