package launchpad

import (
	"o1-launchpad/service/launchpad/api/internal/types"
	"o1-launchpad/service/launchpad/rpc/launchpadclient"
)

func mapLaunchpad(item *launchpadclient.LaunchpadInfo) types.LaunchpadInfo {
	return types.LaunchpadInfo{
		Id: item.Id, ChainId: item.ChainId, Slug: item.Slug, Name: item.Name, Description: item.Description,
		LogoUrl: item.LogoUrl, PrimaryColor: item.PrimaryColor, Owner: item.Owner, Treasury: item.Treasury,
		Active: item.Active, CreatedAt: item.CreatedAt,
	}
}

func mapTokens(items []*launchpadclient.Token) []types.TokenInfo {
	result := make([]types.TokenInfo, 0, len(items))
	for _, item := range items {
		result = append(result, types.TokenInfo{
			ChainId: item.ChainId, Token: item.Token, PoolId: item.PoolId, LaunchpadId: item.LaunchpadId,
			LaunchpadSlug: item.LaunchpadSlug, Creator: item.Creator, Quote: item.Quote, Supply: item.Supply,
			TxHash: item.TxHash, BlockNumber: item.BlockNumber, Status: item.Status, CreatedAt: item.CreatedAt,
		})
	}
	return result
}
