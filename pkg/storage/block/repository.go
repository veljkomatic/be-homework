package block

import (
	"context"

	"github.com/veljkomatic/be-homework/pkg/blockchain"
)

// ReadOnlyBlockRepository is responsible for reading block number
type ReadOnlyBlockRepository interface {
	GetCurrentBlockNumber(ctx context.Context) (blockchain.BlockNumber, error)
}

// WriteBlockRepository is responsible for writing block number
type WriteBlockRepository interface {
	SaveBlockNumber(ctx context.Context, blockNumber blockchain.BlockNumber) error
}

// Repository is responsible for reading and writing block number
type Repository interface {
	ReadOnlyBlockRepository
	WriteBlockRepository
}

var _ Repository = (*repository)(nil)

type repository struct {
	storage Storage
}

func NewRepository(storage Storage) Repository {
	return &repository{
		storage: storage,
	}
}

func (r repository) GetCurrentBlockNumber(ctx context.Context) (blockchain.BlockNumber, error) {
	currentBlockNumber, err := r.storage.Get(ctx)
	if err != nil {
		return blockchain.InvalidBlockNumber, err
	}
	return blockchain.BlockNumber(currentBlockNumber), nil
}

func (r repository) SaveBlockNumber(ctx context.Context, blockNumber blockchain.BlockNumber) error {
	return r.storage.Save(ctx, blockNumber.ToInt64())
}
