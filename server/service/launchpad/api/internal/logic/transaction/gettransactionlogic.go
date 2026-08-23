package transaction

import (
	"context"

	"o1-launchpad/service/launchpad/api/internal/svc"
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTransactionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionLogic {
	return &GetTransactionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTransactionLogic) GetTransaction(req *types.TransactionPath) (resp *types.TransactionInfo, err error) {
	item, err := l.svcCtx.Launchpad.GetTransaction(l.ctx, &launchpadclient.GetTransactionRequest{
		ChainId: req.ChainId, TxHash: req.TxHash,
	})
	if err != nil {
		return nil, err
	}
	return &types.TransactionInfo{
		ChainId: item.ChainId, TxHash: item.TxHash, Kind: item.Kind, Status: item.Status,
		BlockNumber: item.BlockNumber, FailureReason: item.FailureReason, UpdatedAt: item.UpdatedAt,
	}, nil
}
