package logic

import (
	"sync"
	"time"

	"o1-launchpad/service/indexer/internal/svc"
)

type IndexerLogic struct {
	svcCtx *svc.ServiceContext
	stop   chan struct{}
	once   sync.Once
}

func NewIndexerLogic(svcCtx *svc.ServiceContext) *IndexerLogic {
	return &IndexerLogic{svcCtx: svcCtx, stop: make(chan struct{})}
}

func (l *IndexerLogic) Start() {
	ticker := time.NewTicker(l.svcCtx.Config.Indexer.PollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// Chain scanning is added after the generated contract bindings.
		case <-l.stop:
			return
		}
	}
}

func (l *IndexerLogic) Stop() {
	l.once.Do(func() {
		close(l.stop)
		l.svcCtx.Chain.RPC.Close()
		l.svcCtx.DB.Close()
	})
}
