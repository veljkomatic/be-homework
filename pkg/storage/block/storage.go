package block

import (
	"context"
	"sync/atomic"
)

// ReadOnlyStorage is a storage that can only be read from
type ReadOnlyStorage interface {
	Get(ctx context.Context) (int64, error)
}

// WriteStorage is a storage that can be written to
type WriteStorage interface {
	Save(ctx context.Context, lastProcessedBlockNumber int64) error
}

// Storage is a storage that can be read from and written to
type Storage interface {
	ReadOnlyStorage
	WriteStorage
}

var _ Storage = (*inMemoryStorage)(nil)

// inMemoryStorage is a storage that stores last processed block number in memory
type inMemoryStorage struct {
	lastProcessedBlockNumber int64
}

func NewStorage() Storage {
	return &inMemoryStorage{
		lastProcessedBlockNumber: 0,
	}
}

// Get returns last processed block number
func (s *inMemoryStorage) Get(ctx context.Context) (int64, error) {
	return atomic.LoadInt64(&s.lastProcessedBlockNumber), nil
}

// Save saves last processed block number
func (s *inMemoryStorage) Save(ctx context.Context, lastProcessedBlockNumber int64) error {
	atomic.StoreInt64(&s.lastProcessedBlockNumber, lastProcessedBlockNumber)
	return nil
}
