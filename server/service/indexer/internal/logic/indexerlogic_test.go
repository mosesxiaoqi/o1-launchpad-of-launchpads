package logic

import (
	"context"
	"database/sql"
	"math/big"
	"testing"
	"time"

	"o1-launchpad/common/chain"
	commonconfig "o1-launchpad/common/config"
	"o1-launchpad/common/model"
	indexerconfig "o1-launchpad/service/indexer/internal/config"
	"o1-launchpad/service/indexer/internal/svc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type fakeLaunchModel struct {
	checkpoint model.Checkpoint
	hasCP      bool
	saved      model.Checkpoint
	rewind     uint64
}

func (f *fakeLaunchModel) GetCheckpoint(context.Context, int64, []byte) (model.Checkpoint, error) {
	if !f.hasCP {
		return model.Checkpoint{}, sql.ErrNoRows
	}
	return f.checkpoint, nil
}
func (f *fakeLaunchModel) SaveCheckpoint(_ context.Context, checkpoint model.Checkpoint) error {
	f.saved = checkpoint
	return nil
}
func (f *fakeLaunchModel) ApplyLaunchEvent(context.Context, model.LaunchEvent, model.Checkpoint) error {
	return nil
}
func (f *fakeLaunchModel) DeleteLaunchEventsFromBlock(_ context.Context, _ int64, _ []byte, block uint64) error {
	f.rewind = block
	return nil
}

type fakeLogClient struct {
	latest uint64
	query  ethereum.FilterQuery
}

func (f *fakeLogClient) BlockNumber(context.Context) (uint64, error) { return f.latest, nil }
func (f *fakeLogClient) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	return &types.Header{Number: number, Time: 1_700_000_000, Extra: []byte{byte(number.Uint64())}}, nil
}
func (f *fakeLogClient) FilterLogs(_ context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	f.query = query
	return nil, nil
}

func TestRunOnceBoundsBatchAndStopsAtConfirmedHead(t *testing.T) {
	modelFake := &fakeLaunchModel{}
	rpc := &fakeLogClient{latest: 1000}
	logic := newIndexerLogic(indexerServiceContext(), modelFake, rpc)
	if err := logic.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if rpc.query.FromBlock.Uint64() != 100 || rpc.query.ToBlock.Uint64() != 599 {
		t.Fatalf("query range = %s-%s", rpc.query.FromBlock, rpc.query.ToBlock)
	}
	if modelFake.saved.NextBlock != 600 {
		t.Fatalf("next block = %d", modelFake.saved.NextBlock)
	}
}

func TestRunOnceRewindsAtMostConfiguredLookbackOnHashMismatch(t *testing.T) {
	factory := common.HexToAddress("0x1000000000000000000000000000000000000001")
	modelFake := &fakeLaunchModel{hasCP: true, checkpoint: model.Checkpoint{
		ChainID: 84532, Factory: factory.Bytes(), NextBlock: 200, LastBlockHash: common.HexToHash("0x01").Bytes(),
	}}
	rpc := &fakeLogClient{latest: 300}
	logic := newIndexerLogic(indexerServiceContext(), modelFake, rpc)
	if err := logic.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if modelFake.rewind != 136 || rpc.query.FromBlock.Uint64() != 136 {
		t.Fatalf("rewind = %d, query from = %s", modelFake.rewind, rpc.query.FromBlock)
	}
}

func indexerServiceContext() *svc.ServiceContext {
	return &svc.ServiceContext{
		Config: indexerconfig.Config{
			Chain:   commonconfig.Chain{ChainId: 84532, Confirmations: 2},
			Indexer: commonconfig.Indexer{BatchSize: 500, PollInterval: time.Second, ReorgLookback: 64},
		},
		Deployment: chain.Deployment{
			ChainID: 84532, Factory: "0x1000000000000000000000000000000000000001", DeploymentBlock: 100,
		},
	}
}
