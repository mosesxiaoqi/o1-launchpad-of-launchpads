package chain

import (
	"math/big"
	"testing"

	"o1-launchpad/common/chain/bindings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestParseTrustedLaunchedLog(t *testing.T) {
	factory := common.HexToAddress("0x1000000000000000000000000000000000000001")
	token := common.HexToAddress("0x2000000000000000000000000000000000000002")
	creator := common.HexToAddress("0x3000000000000000000000000000000000000003")
	quote := common.HexToAddress("0x4000000000000000000000000000000000000004")
	poolID := common.HexToHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	launchpadID := common.HexToHash("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	txHash := common.HexToHash("0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc")
	parsedABI, err := bindings.MultiTenantLaunchpadFactoryMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	data, err := parsedABI.Events["Launched"].Inputs.NonIndexed().Pack(creator, quote, big.NewInt(1_000_000), big.NewInt(60))
	if err != nil {
		t.Fatal(err)
	}
	log := types.Log{
		Address: factory, Topics: []common.Hash{
			parsedABI.Events["Launched"].ID, common.BytesToHash(token.Bytes()), poolID, launchpadID,
		}, Data: data, BlockNumber: 120, BlockHash: common.HexToHash("0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"),
		TxHash: txHash, Index: 7,
	}
	event, err := ParseLaunchedLog(log, 84532, factory)
	if err != nil {
		t.Fatal(err)
	}
	if common.BytesToAddress(event.Token) != token || common.BytesToHash(event.PoolID) != poolID ||
		common.BytesToHash(event.LaunchpadID) != launchpadID || common.BytesToAddress(event.Creator) != creator ||
		common.BytesToAddress(event.Quote) != quote || event.Supply != "1000000" || event.BlockNumber != 120 ||
		common.BytesToHash(event.TxHash) != txHash || event.LogIndex != 7 || common.BytesToAddress(event.Factory) != factory {
		t.Fatalf("parsed event mismatch: %#v", event)
	}
	log.Address = quote
	if _, err := ParseLaunchedLog(log, 84532, factory); err == nil {
		t.Fatal("expected untrusted factory error")
	}
}
