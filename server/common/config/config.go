package config

import "time"

type Database struct {
	DataSource   string
	MaxOpenConns int
	MaxIdleConns int
}

type Chain struct {
	ChainId       int64
	HttpRpc       string
	WsRpc         string
	Confirmations uint64
}

type Indexer struct {
	BatchSize     uint64
	PollInterval  time.Duration
	ReorgLookback uint64
}
