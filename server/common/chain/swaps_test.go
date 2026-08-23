package chain

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestEncodeUniversalRouterV2ExactInput(t *testing.T) {
	token := common.HexToAddress("0x2000000000000000000000000000000000000002")
	quote := common.HexToAddress("0x1000000000000000000000000000000000000001")
	hook := common.HexToAddress("0x3000000000000000000000000000000000002acc")
	hookData, err := EncodeHookData(common.Address{})
	if err != nil {
		t.Fatal(err)
	}
	data, err := EncodeUniversalRouterExactInput(
		PoolKey(token, quote, hook, 60), true, big.NewInt(1_000), big.NewInt(900), hookData, 1_700_000_000,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 4 || common.Bytes2Hex(data[:4]) != "3593564c" {
		t.Fatalf("execute selector = %x", data[:4])
	}
}
