package logic

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"o1-launchpad/common/chain/bindings"
	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type PrepareLaunchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPrepareLaunchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareLaunchLogic {
	return &PrepareLaunchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PrepareLaunchLogic) PrepareLaunch(in *launchpad.PrepareLaunchRequest) (*launchpad.PreparedTransaction, error) {
	if in.ChainId != 84532 || !common.IsHexAddress(in.Wallet) || !common.IsHexAddress(l.svcCtx.Deployment.Factory) {
		return nil, errors.New("invalid chain, wallet, or factory")
	}
	launchpadID := common.HexToHash(in.LaunchpadId)
	salt := common.HexToHash(in.Salt)
	factoryAddress := common.HexToAddress(l.svcCtx.Deployment.Factory)
	blockNumber, err := l.svcCtx.Chain.RPC.BlockNumber(l.ctx)
	if err != nil {
		return nil, err
	}
	caller, err := bindings.NewMultiTenantLaunchpadFactoryCaller(factoryAddress, l.svcCtx.Chain.RPC)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{Context: l.ctx, BlockNumber: new(big.Int).SetUint64(blockNumber)}
	quote, err := caller.Quote(callOpts)
	if err != nil {
		return nil, err
	}
	version, err := caller.ConfigVersion(callOpts)
	if err != nil {
		return nil, err
	}
	header, err := l.svcCtx.Chain.RPC.HeaderByNumber(l.ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return nil, err
	}
	deadline := header.Time + 30*60
	contractABI, err := bindings.MultiTenantLaunchpadFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	data, err := contractABI.Pack("launch", bindings.MultiTenantLaunchpadFactoryLaunchParams{
		LaunchpadId: launchpadID, Name: in.Name, Symbol: in.Symbol, ContractURI: in.ContractUri,
		Salt: salt, Quote: quote, ExpectedConfigVersion: version, Deadline: deadline,
	})
	if err != nil {
		return nil, err
	}
	return &launchpad.PreparedTransaction{
		ChainId: 84532, From: common.HexToAddress(in.Wallet).Hex(), To: factoryAddress.Hex(),
		Data: "0x" + hex.EncodeToString(data), Value: "0", Deadline: deadline,
		Review: strings.Join([]string{"fixed supply", "permanent liquidity", "1% protocol fee", "0.5% LaaS fee", "16s anti-snipe"}, ", "),
	}, nil
}
