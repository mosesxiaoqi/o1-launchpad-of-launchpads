package chain

import (
	"context"
	"math/big"

	"o1-launchpad/common/chain/bindings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ReadOwedAtLatestBlock(
	ctx context.Context,
	client *ethclient.Client,
	escrow, recipient, currency common.Address,
) (*big.Int, uint64, error) {
	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, 0, err
	}
	contract, err := bindings.NewFeeEscrow(escrow, client)
	if err != nil {
		return nil, 0, err
	}
	amount, err := contract.Owed(&bind.CallOpts{
		Context: ctx, BlockNumber: new(big.Int).SetUint64(blockNumber),
	}, recipient, currency)
	return amount, blockNumber, err
}
