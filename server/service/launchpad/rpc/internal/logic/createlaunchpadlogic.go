package logic

import (
	"context"
	"strings"
	"time"

	"o1-launchpad/common/branding"
	"o1-launchpad/common/chain"
	"o1-launchpad/common/model"
	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateLaunchpadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLaunchpadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLaunchpadLogic {
	return &CreateLaunchpadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLaunchpadLogic) CreateLaunchpad(in *launchpad.CreateLaunchpadRequest) (*launchpad.LaunchpadInfo, error) {
	if err := branding.Validate(in.Slug, in.Name, in.Description, in.LogoUrl, in.PrimaryColor); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if in.ChainId != 84532 || !common.IsHexAddress(in.Wallet) || len(strings.TrimPrefix(in.RegistryTxHash, "0x")) != 64 {
		return nil, status.Error(codes.InvalidArgument, "invalid chain, wallet, or transaction hash")
	}
	if !common.IsHexAddress(l.svcCtx.Deployment.Registry) || l.svcCtx.Chain.RPC == nil {
		return nil, status.Error(codes.FailedPrecondition, "registry deployment is unavailable")
	}
	owner := common.HexToAddress(in.Wallet)
	verified, err := chain.VerifyRegistryCreation(
		l.ctx, l.svcCtx.Chain.RPC, in.ChainId, common.HexToAddress(l.svcCtx.Deployment.Registry),
		common.HexToHash(in.RegistryTxHash), in.Slug, owner, l.svcCtx.Config.Chain.Confirmations,
	)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	createdAt := time.Now().UTC()
	record := model.Launchpad{
		ChainID: in.ChainId, LaunchpadID: verified.LaunchpadID.Bytes(), Slug: in.Slug, Name: in.Name,
		Description: in.Description, LogoURL: in.LogoUrl, PrimaryColor: in.PrimaryColor,
		Owner: verified.Owner.Bytes(), Treasury: verified.Treasury.Bytes(), RegistryTxHash: common.HexToHash(in.RegistryTxHash).Bytes(),
		RegistryBlockNumber: verified.BlockNumber, Active: true, CreatedAt: createdAt,
	}
	if err := l.svcCtx.Model.CreateLaunchpad(l.ctx, record); err != nil {
		return nil, status.Error(codes.AlreadyExists, "launchpad already exists")
	}
	return &launchpad.LaunchpadInfo{
		Id: verified.LaunchpadID.Hex(), ChainId: in.ChainId, Slug: in.Slug, Name: in.Name,
		Description: in.Description, LogoUrl: in.LogoUrl, PrimaryColor: in.PrimaryColor,
		Owner: verified.Owner.Hex(), Treasury: verified.Treasury.Hex(), Active: true, CreatedAt: createdAt.Format(time.RFC3339),
	}, nil
}
