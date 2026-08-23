package logic

import (
	"context"
	"encoding/hex"
	"errors"

	"o1-launchpad/common/chain/bindings"
	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type PrepareFeeClaimLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPrepareFeeClaimLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareFeeClaimLogic {
	return &PrepareFeeClaimLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PrepareFeeClaimLogic) PrepareFeeClaim(in *launchpad.PrepareFeeClaimRequest) (*launchpad.PreparedTransaction, error) {
	if in.ChainId != 84532 || !common.IsHexAddress(in.Wallet) || !common.IsHexAddress(in.Recipient) ||
		!common.IsHexAddress(in.Currency) || !common.IsHexAddress(l.svcCtx.Deployment.FeeEscrow) {
		return nil, errors.New("invalid claim request")
	}
	wallet := common.HexToAddress(in.Wallet)
	recipient := common.HexToAddress(in.Recipient)
	if wallet != recipient {
		return nil, errors.New("wallet must claim its own balance")
	}
	contractABI, err := bindings.FeeEscrowMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	data, err := contractABI.Pack("claim", recipient, common.HexToAddress(in.Currency))
	if err != nil {
		return nil, err
	}
	return &launchpad.PreparedTransaction{
		ChainId: 84532, From: wallet.Hex(), To: common.HexToAddress(l.svcCtx.Deployment.FeeEscrow).Hex(),
		Data: "0x" + hex.EncodeToString(data), Value: "0", Review: "claim accrued fee balance",
	}, nil
}
