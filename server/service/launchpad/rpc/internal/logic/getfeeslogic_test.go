package logic

import (
	"context"
	"testing"

	commonconfig "o1-launchpad/common/config"
	"o1-launchpad/service/launchpad/rpc/internal/config"
	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"
)

func TestGetFeesRejectsInvalidAddressBeforeChainCall(t *testing.T) {
	logic := NewGetFeesLogic(context.Background(), &svc.ServiceContext{
		Config: config.Config{Chain: commonconfig.Chain{ChainId: 84532}},
	})
	_, err := logic.GetFees(&launchpad.GetFeesRequest{
		ChainId: 84532, Recipient: "invalid", Currency: "0x0000000000000000000000000000000000000000",
	})
	if err == nil {
		t.Fatal("expected invalid fee query error")
	}
}
