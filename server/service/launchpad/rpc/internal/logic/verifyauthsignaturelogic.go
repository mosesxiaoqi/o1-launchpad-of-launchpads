package logic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	commonauth "o1-launchpad/common/auth"
	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VerifyAuthSignatureLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyAuthSignatureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyAuthSignatureLogic {
	return &VerifyAuthSignatureLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyAuthSignatureLogic) VerifyAuthSignature(in *launchpad.VerifyAuthSignatureRequest) (*launchpad.VerifyAuthSignatureResponse, error) {
	if !common.IsHexAddress(in.Address) {
		return nil, status.Error(codes.InvalidArgument, "invalid wallet address")
	}
	nonce, err := l.svcCtx.Model.GetAuthNonce(l.ctx, in.ChallengeId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.Unauthenticated, "challenge missing or already used")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "load auth challenge")
	}
	challenge := commonauth.Challenge{
		Domain: nonce.Domain, URI: nonce.URI, ChainID: nonce.ChainID, Address: common.BytesToAddress(nonce.Wallet),
		Nonce: nonce.Nonce, Message: nonce.Message, IssuedAt: nonce.IssuedAt, ExpiresAt: nonce.ExpiresAt,
	}
	now := time.Now().UTC()
	address := common.HexToAddress(in.Address)
	if err := commonauth.Verify(challenge, in.Signature, address, now); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if _, err := l.svcCtx.Model.ConsumeAuthNonce(l.ctx, in.ChallengeId, now); err != nil {
		return nil, status.Error(codes.Unauthenticated, "challenge already used")
	}
	sessionExpiry := now.Add(24 * time.Hour)
	return &launchpad.VerifyAuthSignatureResponse{
		SessionSubject: "84532:" + address.Hex(), ChainId: 84532, Address: address.Hex(),
		ExpiresAt: sessionExpiry.Format(time.RFC3339),
	}, nil
}
