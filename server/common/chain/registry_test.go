package chain

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type fakeReceiptReader struct {
	receipt *types.Receipt
	block   uint64
}

func (f fakeReceiptReader) TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error) {
	return f.receipt, nil
}

func (f fakeReceiptReader) BlockNumber(context.Context) (uint64, error) {
	return f.block, nil
}

func TestVerifyRegistryCreationRejectsMismatchesAndLowConfirmations(t *testing.T) {
	registry := common.HexToAddress("0x1000000000000000000000000000000000000001")
	owner := common.HexToAddress("0x2000000000000000000000000000000000000002")
	treasury := common.HexToAddress("0x3000000000000000000000000000000000000003")
	txHash := common.HexToHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	receipt := registryReceipt(t, registry, owner, treasury, "ai-pad", 100)

	if _, err := VerifyRegistryCreation(
		context.Background(), fakeReceiptReader{receipt: receipt, block: 101}, 84532, registry, txHash, "ai-pad", owner, 2,
	); err != nil {
		t.Fatalf("VerifyRegistryCreation() error = %v", err)
	}

	tests := []struct {
		name     string
		chainID  int64
		registry common.Address
		slug     string
		owner    common.Address
		block    uint64
	}{
		{name: "wrong chain", chainID: 1, registry: registry, slug: "ai-pad", owner: owner, block: 101},
		{name: "wrong registry", chainID: 84532, registry: common.HexToAddress("0x4000000000000000000000000000000000000004"), slug: "ai-pad", owner: owner, block: 101},
		{name: "wrong slug", chainID: 84532, registry: registry, slug: "other-pad", owner: owner, block: 101},
		{name: "wrong owner", chainID: 84532, registry: registry, slug: "ai-pad", owner: treasury, block: 101},
		{name: "one confirmation", chainID: 84532, registry: registry, slug: "ai-pad", owner: owner, block: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := VerifyRegistryCreation(
				context.Background(), fakeReceiptReader{receipt: receipt, block: tt.block}, tt.chainID,
				tt.registry, txHash, tt.slug, tt.owner, 2,
			); err == nil {
				t.Fatal("expected verification error")
			}
		})
	}

	reverted := *receipt
	reverted.Status = types.ReceiptStatusFailed
	if _, err := VerifyRegistryCreation(
		context.Background(), fakeReceiptReader{receipt: &reverted, block: 101}, 84532,
		registry, txHash, "ai-pad", owner, 2,
	); err == nil {
		t.Fatal("expected reverted receipt error")
	}
}

func registryReceipt(
	t *testing.T, registry, owner, treasury common.Address, slug string, block uint64,
) *types.Receipt {
	t.Helper()
	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	data, err := (abi.Arguments{{Type: stringType}}).Pack(slug)
	if err != nil {
		t.Fatal(err)
	}
	id, err := DeriveLaunchpadID(84532, slug)
	if err != nil {
		t.Fatal(err)
	}
	return &types.Receipt{
		Status: types.ReceiptStatusSuccessful, BlockNumber: new(big.Int).SetUint64(block),
		Logs: []*types.Log{{
			Address: registry,
			Topics: []common.Hash{
				crypto.Keccak256Hash([]byte("LaunchpadCreated(bytes32,address,address,string)")),
				id, common.BytesToHash(owner.Bytes()), common.BytesToHash(treasury.Bytes()),
			},
			Data: data,
		}},
	}
}
