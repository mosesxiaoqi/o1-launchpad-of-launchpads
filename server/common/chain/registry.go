package chain

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type ReceiptReader interface {
	TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error)
	BlockNumber(context.Context) (uint64, error)
}

type RegistryCreation struct {
	LaunchpadID common.Hash
	Owner       common.Address
	Treasury    common.Address
	Slug        string
	BlockNumber uint64
}

var launchpadCreatedTopic = crypto.Keccak256Hash([]byte("LaunchpadCreated(bytes32,address,address,string)"))

func DeriveLaunchpadID(chainID int64, slug string) (common.Hash, error) {
	uintType, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return common.Hash{}, err
	}
	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return common.Hash{}, err
	}
	encoded, err := (abi.Arguments{{Type: uintType}, {Type: stringType}}).Pack(big.NewInt(chainID), slug)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(encoded), nil
}

func VerifyRegistryCreation(
	ctx context.Context,
	reader ReceiptReader,
	chainID int64,
	registry common.Address,
	txHash common.Hash,
	slug string,
	owner common.Address,
	requiredConfirmations uint64,
) (RegistryCreation, error) {
	if chainID != 84532 || registry == (common.Address{}) || owner == (common.Address{}) {
		return RegistryCreation{}, errors.New("invalid chain, registry, or owner")
	}
	receipt, err := reader.TransactionReceipt(ctx, txHash)
	if err != nil {
		return RegistryCreation{}, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful || receipt.BlockNumber == nil {
		return RegistryCreation{}, errors.New("registry transaction reverted")
	}
	currentBlock, err := reader.BlockNumber(ctx)
	if err != nil {
		return RegistryCreation{}, err
	}
	blockNumber := receipt.BlockNumber.Uint64()
	if currentBlock < blockNumber || currentBlock-blockNumber+1 < requiredConfirmations {
		return RegistryCreation{}, errors.New("insufficient confirmations")
	}
	expectedID, err := DeriveLaunchpadID(chainID, slug)
	if err != nil {
		return RegistryCreation{}, err
	}
	stringType, _ := abi.NewType("string", "", nil)
	arguments := abi.Arguments{{Type: stringType}}
	for _, event := range receipt.Logs {
		if event.Address != registry || len(event.Topics) != 4 || event.Topics[0] != launchpadCreatedTopic {
			continue
		}
		values, err := arguments.Unpack(event.Data)
		if err != nil || len(values) != 1 {
			continue
		}
		eventSlug, ok := values[0].(string)
		eventOwner := common.BytesToAddress(event.Topics[2].Bytes())
		if !ok || event.Topics[1] != expectedID || eventOwner != owner || eventSlug != slug {
			continue
		}
		return RegistryCreation{
			LaunchpadID: event.Topics[1], Owner: eventOwner,
			Treasury: common.BytesToAddress(event.Topics[3].Bytes()), Slug: eventSlug, BlockNumber: blockNumber,
		}, nil
	}
	return RegistryCreation{}, fmt.Errorf("matching LaunchpadCreated event not found")
}
