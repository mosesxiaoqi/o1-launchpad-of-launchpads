package swap

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type Quote struct {
	ID            string
	Wallet        common.Address
	Token         common.Address
	QuoteCurrency common.Address
	Hook          common.Address
	ZeroForOne    bool
	AmountIn      *big.Int
	AmountOut     *big.Int
	HookData      []byte
	ExpiresAt     time.Time
}

type Store struct {
	mu     sync.RWMutex
	quotes map[string]Quote
}

func NewStore() *Store {
	return &Store{quotes: make(map[string]Quote)}
}

func (s *Store) Put(quote Quote) Quote {
	s.mu.Lock()
	defer s.mu.Unlock()
	quote.ID = uuid.NewString()
	s.quotes[quote.ID] = quote
	return quote
}

func (s *Store) Get(id string, now time.Time) (Quote, error) {
	s.mu.RLock()
	quote, ok := s.quotes[id]
	s.mu.RUnlock()
	if !ok || !now.Before(quote.ExpiresAt) {
		return Quote{}, errors.New("quote missing or expired")
	}
	return quote, nil
}
