package logic

import (
	"context"
	"time"

	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionLogic {
	return &GetTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTransactionLogic) GetTransaction(in *launchpad.GetTransactionRequest) (*launchpad.Transaction, error) {
	item, err := l.svcCtx.Model.GetTransaction(l.ctx, in.ChainId, common.HexToHash(in.TxHash).Bytes())
	if err != nil {
		return nil, err
	}
	return &launchpad.Transaction{
		ChainId: item.ChainID, TxHash: common.BytesToHash(item.TxHash).Hex(), Kind: item.Kind, Status: item.Status,
		BlockNumber: item.BlockNumber, FailureReason: item.FailureReason, UpdatedAt: item.UpdatedAt.Format(time.RFC3339),
	}, nil
}
