package logic

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"time"

	"o1-launchpad/common/chain"
	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type PrepareSwapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPrepareSwapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareSwapLogic {
	return &PrepareSwapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PrepareSwapLogic) PrepareSwap(in *launchpad.PrepareSwapRequest) (*launchpad.PreparedTransaction, error) {
	if in.ChainId != 84532 || !common.IsHexAddress(in.Wallet) || in.SlippageBps > 500 {
		return nil, errors.New("invalid swap preparation request")
	}
	quote, err := l.svcCtx.Quotes.Get(in.QuoteId, time.Now().UTC())
	if err != nil || quote.Wallet != common.HexToAddress(in.Wallet) {
		return nil, errors.New("quote missing, expired, or belongs to another wallet")
	}
	minimumOut := new(big.Int).Mul(quote.AmountOut, new(big.Int).SetUint64(uint64(10_000-in.SlippageBps)))
	minimumOut.Div(minimumOut, big.NewInt(10_000))
	deadline := uint64(time.Now().Add(30 * time.Minute).Unix())
	poolKey := chain.PoolKey(quote.Token, quote.QuoteCurrency, quote.Hook, l.svcCtx.Deployment.TickSpacing)
	data, err := chain.EncodeUniversalRouterExactInput(
		poolKey, quote.ZeroForOne, quote.AmountIn, minimumOut, quote.HookData, deadline,
	)
	if err != nil {
		return nil, err
	}
	value := "0"
	inputCurrency := poolKey.Currency1
	if quote.ZeroForOne {
		inputCurrency = poolKey.Currency0
	}
	if inputCurrency == (common.Address{}) {
		value = quote.AmountIn.String()
	}
	return &launchpad.PreparedTransaction{
		ChainId: 84532, From: quote.Wallet.Hex(), To: common.HexToAddress(l.svcCtx.Deployment.UniversalRouter).Hex(),
		Data: "0x" + hex.EncodeToString(data), Value: value, Deadline: deadline,
		Review: "exact-input Uniswap v4 swap; max slippage " + new(big.Int).SetUint64(uint64(in.SlippageBps)).String() + " bps",
	}, nil
}
