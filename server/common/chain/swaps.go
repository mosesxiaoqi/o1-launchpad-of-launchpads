package chain

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type V4PoolKey struct {
	Currency0   common.Address
	Currency1   common.Address
	Fee         *big.Int
	TickSpacing *big.Int
	Hooks       common.Address
}

type quoteExactSingleParams struct {
	PoolKey     V4PoolKey
	ZeroForOne  bool
	ExactAmount *big.Int
	HookData    []byte
}

type exactInputSingleParams struct {
	PoolKey          V4PoolKey
	ZeroForOne       bool
	AmountIn         *big.Int
	AmountOutMinimum *big.Int
	HookData         []byte
}

const quoterABI = `[{"type":"function","name":"quoteExactInputSingle","stateMutability":"nonpayable","inputs":[{"name":"params","type":"tuple","components":[{"name":"poolKey","type":"tuple","components":[{"name":"currency0","type":"address"},{"name":"currency1","type":"address"},{"name":"fee","type":"uint24"},{"name":"tickSpacing","type":"int24"},{"name":"hooks","type":"address"}]},{"name":"zeroForOne","type":"bool"},{"name":"exactAmount","type":"uint128"},{"name":"hookData","type":"bytes"}]}],"outputs":[{"name":"amountOut","type":"uint256"},{"name":"gasEstimate","type":"uint256"}]}]`

const universalRouterABI = `[{"type":"function","name":"execute","stateMutability":"payable","inputs":[{"name":"commands","type":"bytes"},{"name":"inputs","type":"bytes[]"},{"name":"deadline","type":"uint256"}],"outputs":[]}]`

func PoolKey(token, quote, hook common.Address, tickSpacing int64) V4PoolKey {
	if bytesLess(token.Bytes(), quote.Bytes()) {
		return V4PoolKey{Currency0: token, Currency1: quote, Fee: big.NewInt(0), TickSpacing: big.NewInt(tickSpacing), Hooks: hook}
	}
	return V4PoolKey{Currency0: quote, Currency1: token, Fee: big.NewInt(0), TickSpacing: big.NewInt(tickSpacing), Hooks: hook}
}

func EncodeHookData(referrer common.Address) ([]byte, error) {
	addressType, _ := abi.NewType("address", "", nil)
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	return (abi.Arguments{{Type: addressType}, {Type: bytes32Type}}).Pack(referrer, [32]byte{})
}

func QuoteExactInputSingle(
	ctx context.Context,
	client *ethclient.Client,
	quoter, from common.Address,
	poolKey V4PoolKey,
	zeroForOne bool,
	amountIn *big.Int,
	hookData []byte,
	blockNumber *big.Int,
) (*big.Int, error) {
	if amountIn.Sign() <= 0 || amountIn.BitLen() > 128 {
		return nil, errors.New("amount must fit uint128")
	}
	contractABI, err := abi.JSON(strings.NewReader(quoterABI))
	if err != nil {
		return nil, err
	}
	data, err := contractABI.Pack("quoteExactInputSingle", quoteExactSingleParams{
		PoolKey: poolKey, ZeroForOne: zeroForOne, ExactAmount: amountIn, HookData: hookData,
	})
	if err != nil {
		return nil, err
	}
	result, err := client.CallContract(ctx, ethereum.CallMsg{From: from, To: &quoter, Data: data}, blockNumber)
	if err != nil {
		return nil, err
	}
	values, err := contractABI.Unpack("quoteExactInputSingle", result)
	if err != nil || len(values) != 2 {
		return nil, errors.New("invalid quoter response")
	}
	amountOut, ok := values[0].(*big.Int)
	if !ok {
		return nil, errors.New("invalid quote amount")
	}
	return amountOut, nil
}

func EncodeUniversalRouterExactInput(
	poolKey V4PoolKey,
	zeroForOne bool,
	amountIn, amountOutMinimum *big.Int,
	hookData []byte,
	deadline uint64,
) ([]byte, error) {
	poolComponents := []abi.ArgumentMarshaling{
		{Name: "currency0", Type: "address"}, {Name: "currency1", Type: "address"}, {Name: "fee", Type: "uint24"},
		{Name: "tickSpacing", Type: "int24"}, {Name: "hooks", Type: "address"},
	}
	exactType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "poolKey", Type: "tuple", Components: poolComponents}, {Name: "zeroForOne", Type: "bool"},
		{Name: "amountIn", Type: "uint128"}, {Name: "amountOutMinimum", Type: "uint128"}, {Name: "hookData", Type: "bytes"},
	})
	if err != nil {
		return nil, err
	}
	exact, err := (abi.Arguments{{Type: exactType}}).Pack(exactInputSingleParams{
		PoolKey: poolKey, ZeroForOne: zeroForOne, AmountIn: amountIn, AmountOutMinimum: amountOutMinimum, HookData: hookData,
	})
	if err != nil {
		return nil, err
	}
	addressType, _ := abi.NewType("address", "", nil)
	uintType, _ := abi.NewType("uint256", "", nil)
	inputCurrency := poolKey.Currency1
	outputCurrency := poolKey.Currency0
	if zeroForOne {
		inputCurrency, outputCurrency = poolKey.Currency0, poolKey.Currency1
	}
	settle, err := (abi.Arguments{{Type: addressType}, {Type: uintType}}).Pack(inputCurrency, amountIn)
	if err != nil {
		return nil, err
	}
	take, err := (abi.Arguments{{Type: addressType}, {Type: uintType}}).Pack(outputCurrency, amountOutMinimum)
	if err != nil {
		return nil, err
	}
	bytesType, _ := abi.NewType("bytes", "", nil)
	bytesArrayType, _ := abi.NewType("bytes[]", "", nil)
	v4Input, err := (abi.Arguments{{Type: bytesType}, {Type: bytesArrayType}}).Pack(
		[]byte{0x06, 0x0c, 0x0f}, [][]byte{exact, settle, take},
	)
	if err != nil {
		return nil, err
	}
	routerABI, err := abi.JSON(strings.NewReader(universalRouterABI))
	if err != nil {
		return nil, err
	}
	return routerABI.Pack("execute", []byte{0x10}, [][]byte{v4Input}, new(big.Int).SetUint64(deadline))
}

func bytesLess(a, b []byte) bool {
	for i := range a {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}
