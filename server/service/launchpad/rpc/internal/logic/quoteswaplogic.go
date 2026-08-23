package logic

import (
	"context"
	"errors"
	"math/big"
	"time"

	"o1-launchpad/common/chain"
	"o1-launchpad/common/swap"
	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type QuoteSwapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuoteSwapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuoteSwapLogic {
	return &QuoteSwapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QuoteSwapLogic) QuoteSwap(in *launchpad.QuoteSwapRequest) (*launchpad.QuoteSwapResponse, error) {
	amountIn, ok := new(big.Int).SetString(in.Amount, 10)
	if !ok || amountIn.Sign() <= 0 || in.ChainId != 84532 || !common.IsHexAddress(in.Token) || !common.IsHexAddress(in.Wallet) {
		return nil, errors.New("invalid quote request")
	}
	item, err := l.svcCtx.Model.GetToken(l.ctx, in.ChainId, common.HexToAddress(in.Token).Bytes())
	if err != nil {
		return nil, err
	}
	token := common.BytesToAddress(item.Token)
	quoteCurrency := common.BytesToAddress(item.Quote)
	hook := common.HexToAddress(l.svcCtx.Deployment.Hook)
	poolKey := chain.PoolKey(token, quoteCurrency, hook, l.svcCtx.Deployment.TickSpacing)
	inputCurrency := token
	if in.Buy {
		inputCurrency = quoteCurrency
	}
	zeroForOne := inputCurrency == poolKey.Currency0
	referrer := common.Address{}
	if in.Referrer != "" {
		if !common.IsHexAddress(in.Referrer) {
			return nil, errors.New("invalid referrer")
		}
		referrer = common.HexToAddress(in.Referrer)
	}
	hookData, err := chain.EncodeHookData(referrer)
	if err != nil {
		return nil, err
	}
	blockNumber, err := l.svcCtx.Chain.RPC.BlockNumber(l.ctx)
	if err != nil {
		return nil, err
	}
	amountOut, err := chain.QuoteExactInputSingle(
		l.ctx, l.svcCtx.Chain.RPC, common.HexToAddress(l.svcCtx.Deployment.Quoter), common.HexToAddress(in.Wallet),
		poolKey, zeroForOne, amountIn, hookData, new(big.Int).SetUint64(blockNumber),
	)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().UTC().Add(30 * time.Second)
	stored := l.svcCtx.Quotes.Put(swap.Quote{
		Wallet: common.HexToAddress(in.Wallet), Token: token, QuoteCurrency: quoteCurrency, Hook: hook,
		ZeroForOne: zeroForOne, AmountIn: amountIn, AmountOut: amountOut, HookData: hookData, ExpiresAt: expiresAt,
	})
	fee := new(big.Int).Mul(amountIn, big.NewInt(150))
	fee.Div(fee, big.NewInt(10_000))
	return &launchpad.QuoteSwapResponse{
		QuoteId: stored.ID, AmountIn: amountIn.String(), AmountOut: amountOut.String(), Fee: fee.String(),
		ExpiresAt: uint64(expiresAt.Unix()),
	}, nil
}
