package logic

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	commonauth "o1-launchpad/common/auth"
	"o1-launchpad/common/model"
	"o1-launchpad/service/launchpad/rpc/internal/svc"
	"o1-launchpad/service/launchpad/rpc/pb/launchpad"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateAuthChallengeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAuthChallengeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAuthChallengeLogic {
	return &CreateAuthChallengeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAuthChallengeLogic) CreateAuthChallenge(in *launchpad.CreateAuthChallengeRequest) (*launchpad.CreateAuthChallengeResponse, error) {
	if !common.IsHexAddress(in.Address) {
		return nil, status.Error(codes.InvalidArgument, "invalid wallet address")
	}
	now := time.Now().UTC()
	challenge, err := commonauth.NewChallenge(
		in.Domain, in.Uri, in.ChainId, common.HexToAddress(in.Address), now, 5*time.Minute,
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	id := uuid.NewString()
	nonceHash := sha256.Sum256([]byte(challenge.Nonce))
	err = l.svcCtx.Model.CreateAuthNonce(l.ctx, model.AuthNonce{
		ID: id, ChainID: challenge.ChainID, Wallet: challenge.Address.Bytes(), NonceHash: nonceHash[:],
		Domain: challenge.Domain, URI: challenge.URI, Nonce: challenge.Nonce, Message: challenge.Message,
		IssuedAt: challenge.IssuedAt, ExpiresAt: challenge.ExpiresAt,
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, status.Error(codes.Canceled, err.Error())
		}
		return nil, status.Error(codes.Internal, "save auth challenge")
	}
	return &launchpad.CreateAuthChallengeResponse{
		ChallengeId: id, Message: challenge.Message, ExpiresAt: challenge.ExpiresAt.Format(time.RFC3339),
	}, nil
}
