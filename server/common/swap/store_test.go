package swap

import (
	"testing"
	"time"
)

func TestStoreRejectsExpiredQuote(t *testing.T) {
	now := time.Now().UTC()
	store := NewStore()
	quote := store.Put(Quote{ExpiresAt: now.Add(time.Second)})
	if _, err := store.Get(quote.ID, now); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(quote.ID, now.Add(time.Second)); err == nil {
		t.Fatal("expected expired quote error")
	}
}
