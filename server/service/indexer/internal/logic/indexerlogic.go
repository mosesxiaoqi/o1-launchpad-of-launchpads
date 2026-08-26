package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"o1-launchpad/common/chain"
	"o1-launchpad/common/chain/bindings"
	"o1-launchpad/common/model"
	"o1-launchpad/service/indexer/internal/svc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type launchModel interface {
	GetCheckpoint(context.Context, int64, []byte) (model.Checkpoint, error)
	SaveCheckpoint(context.Context, model.Checkpoint) error
	ApplyLaunchEvent(context.Context, model.LaunchEvent, model.Checkpoint) error
	DeleteLaunchEventsFromBlock(context.Context, int64, []byte, uint64) error
}

type logClient interface {
	BlockNumber(context.Context) (uint64, error)
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
}

type IndexerLogic struct {
	svcCtx *svc.ServiceContext
	model  launchModel
	rpc    logClient
	stop   chan struct{}
	once   sync.Once
}

func NewIndexerLogic(svcCtx *svc.ServiceContext) *IndexerLogic {
	return newIndexerLogic(svcCtx, svcCtx.Model, svcCtx.Chain.RPC)
}

func newIndexerLogic(svcCtx *svc.ServiceContext, launchModel launchModel, rpc logClient) *IndexerLogic {
	return &IndexerLogic{svcCtx: svcCtx, model: launchModel, rpc: rpc, stop: make(chan struct{})}
}

func (l *IndexerLogic) Start() {
	ticker := time.NewTicker(l.svcCtx.Config.Indexer.PollInterval)
	defer ticker.Stop()
	for {
		if err := l.RunOnce(context.Background()); err != nil {
			fmt.Printf("indexer scan failed: %v\n", err)
		}
		select {
		case <-ticker.C:
		case <-l.stop:
			return
		}
	}
}

func (l *IndexerLogic) RunOnce(ctx context.Context) error {
	config := l.svcCtx.Config
	deployment := l.svcCtx.Deployment
	if !common.IsHexAddress(deployment.Factory) || deployment.DeploymentBlock == 0 {
		return errors.New("factory deployment is not configured")
	}
	factory := common.HexToAddress(deployment.Factory)
	factoryBytes := factory.Bytes()
	latest, err := l.rpc.BlockNumber(ctx)
	if err != nil {
		return err
	}
	if latest <= config.Chain.Confirmations {
		return nil
	}
	target := latest - config.Chain.Confirmations
	checkpoint, err := l.model.GetCheckpoint(ctx, config.Chain.ChainId, factoryBytes)
	if errors.Is(err, sql.ErrNoRows) {
		checkpoint = model.Checkpoint{
			ChainID: config.Chain.ChainId, Factory: factoryBytes, NextBlock: target,
		}
	} else if err != nil {
		return err
	}

	if checkpoint.NextBlock > deployment.DeploymentBlock && len(checkpoint.LastBlockHash) > 0 {
		previous := checkpoint.NextBlock - 1
		header, err := l.rpc.HeaderByNumber(ctx, new(big.Int).SetUint64(previous))
		if err != nil {
			return err
		}
		if header.Hash() != common.BytesToHash(checkpoint.LastBlockHash) {
			rewind := deployment.DeploymentBlock
			if checkpoint.NextBlock > config.Indexer.ReorgLookback+deployment.DeploymentBlock {
				rewind = checkpoint.NextBlock - config.Indexer.ReorgLookback
			}
			if err := l.model.DeleteLaunchEventsFromBlock(ctx, config.Chain.ChainId, factoryBytes, rewind); err != nil {
				return err
			}
			checkpoint.NextBlock = rewind
			checkpoint.LastBlockHash = nil
		}
	}

	if checkpoint.NextBlock > target {
		return nil
	}
	end := checkpoint.NextBlock + config.Indexer.BatchSize - 1
	if end > target {
		end = target
	}
	parsedABI, err := bindings.MultiTenantLaunchpadFactoryMetaData.GetAbi()
	if err != nil {
		return err
	}
	logs, err := l.rpc.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(checkpoint.NextBlock), ToBlock: new(big.Int).SetUint64(end),
		Addresses: []common.Address{factory}, Topics: [][]common.Hash{{parsedABI.Events["Launched"].ID}},
	})
	if err != nil {
		return err
	}
	endHeader, err := l.rpc.HeaderByNumber(ctx, new(big.Int).SetUint64(end))
	if err != nil {
		return err
	}
	nextCheckpoint := model.Checkpoint{
		ChainID: config.Chain.ChainId, Factory: factoryBytes, NextBlock: end + 1, LastBlockHash: endHeader.Hash().Bytes(),
	}
	for _, entry := range logs {
		event, err := chain.ParseLaunchedLog(entry, config.Chain.ChainId, factory)
		if err != nil {
			return err
		}
		header, err := l.rpc.HeaderByNumber(ctx, new(big.Int).SetUint64(entry.BlockNumber))
		if err != nil {
			return err
		}
		event.CreatedAt = time.Unix(int64(header.Time), 0).UTC()
		if err := l.model.ApplyLaunchEvent(ctx, event, nextCheckpoint); err != nil {
			return err
		}
	}
	if len(logs) == 0 {
		return l.model.SaveCheckpoint(ctx, nextCheckpoint)
	}
	return nil
}

func (l *IndexerLogic) Stop() {
	l.once.Do(func() {
		close(l.stop)
		l.svcCtx.Chain.RPC.Close()
		l.svcCtx.DB.Close()
	})
}
