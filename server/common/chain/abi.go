package chain

import (
	"errors"

	"o1-launchpad/common/chain/bindings"
	"o1-launchpad/common/model"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func ParseLaunchedLog(log types.Log, chainID int64, trustedFactory common.Address) (model.LaunchEvent, error) {
	if log.Address != trustedFactory {
		return model.LaunchEvent{}, errors.New("untrusted factory log")
	}
	filterer, err := bindings.NewMultiTenantLaunchpadFactoryFilterer(trustedFactory, nil)
	if err != nil {
		return model.LaunchEvent{}, err
	}
	launched, err := filterer.ParseLaunched(log)
	if err != nil {
		return model.LaunchEvent{}, err
	}
	return model.LaunchEvent{
		ChainID: chainID, LaunchpadID: launched.LaunchpadId[:], Token: launched.Token.Bytes(),
		PoolID: launched.PoolId[:], Creator: launched.Creator.Bytes(), Quote: launched.Quote.Bytes(),
		Factory: trustedFactory.Bytes(), Supply: launched.Supply.String(), TxHash: log.TxHash.Bytes(),
		BlockNumber: log.BlockNumber, BlockHash: log.BlockHash.Bytes(), LogIndex: uint32(log.Index), Status: "confirmed",
	}, nil
}
