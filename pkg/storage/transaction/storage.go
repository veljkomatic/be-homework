package transaction

import (
	"context"
	"github.com/veljkomatic/be-homework/pkg/blockchain"
	"sync"
)

// ReadOnlyStorage is responsible for reading transactions
type ReadOnlyStorage interface {
	Get(ctx context.Context, key string) ([]*blockchain.Transaction, error)
}

// WriteStorage is responsible for writing transactions
// TODO in future do not use blockchain.Transaction, but some model representation of transaction
type WriteStorage interface {
	InsertBatch(ctx context.Context, data map[string][]*blockchain.Transaction) error
}

// Storage is responsible for reading and writing transactions
type Storage interface {
	ReadOnlyStorage
	WriteStorage
}

var _ Storage = (*inMemoryStorage)(nil)

type inMemoryStorage struct {
	transactions map[string][]*blockchain.Transaction
	mutex        sync.RWMutex
}

func NewStorage() Storage {
	return &inMemoryStorage{
		transactions: make(map[string][]*blockchain.Transaction),
	}
}

func (s *inMemoryStorage) Get(ctx context.Context, key string) ([]*blockchain.Transaction, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	transactions, ok := s.transactions[key]
	if !ok {
		return nil, nil
	}
	return transactions, nil
}

func (s *inMemoryStorage) InsertBatch(ctx context.Context, data map[string][]*blockchain.Transaction) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for key, value := range data {
		s.transactions[key] = append(s.transactions[key], value...)
	}
	return nil
}
