package logic

import (
	"time"

	"o1-launchpad/common/model"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/common"
)

func launchpadInfo(item model.Launchpad) *launchpad.LaunchpadInfo {
	return &launchpad.LaunchpadInfo{
		Id: common.BytesToHash(item.LaunchpadID).Hex(), ChainId: item.ChainID, Slug: item.Slug, Name: item.Name,
		Description: item.Description, LogoUrl: item.LogoURL, PrimaryColor: item.PrimaryColor,
		Owner: common.BytesToAddress(item.Owner).Hex(), Treasury: common.BytesToAddress(item.Treasury).Hex(),
		Active: item.Active, CreatedAt: item.CreatedAt.Format(time.RFC3339),
	}
}

func tokenInfo(item model.TokenLaunch) *launchpad.Token {
	return &launchpad.Token{
		ChainId: item.ChainID, Token: common.BytesToAddress(item.Token).Hex(), PoolId: common.BytesToHash(item.PoolID).Hex(),
		LaunchpadId: common.BytesToHash(item.LaunchpadID).Hex(), LaunchpadSlug: item.LaunchpadSlug,
		Creator: common.BytesToAddress(item.Creator).Hex(), Quote: common.BytesToAddress(item.Quote).Hex(),
		Supply: item.Supply, TxHash: common.BytesToHash(item.TxHash).Hex(), BlockNumber: item.BlockNumber,
		Status: item.Status, CreatedAt: item.CreatedAt.Format(time.RFC3339),
	}
}
