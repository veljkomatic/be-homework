package transaction

import (
	"context"
	"fmt"

	"github.com/veljkomatic/be-homework/pkg/blockchain"
)

// ReadOnlyRepository is responsible for reading transactions
type ReadOnlyRepository interface {
	// GetTransactions list of inbound or outbound transactions for an address
	// TODO Add pagination and filtering in the future
	GetTransactions(ctx context.Context, address string) ([]*blockchain.Transaction, error)
}

// WriteRepository is responsible for writing transactions
type WriteRepository interface {
	InsertTransactions(ctx context.Context, transactions []*AddressTransaction) error
}

// Repository is responsible for reading and writing transactions
type Repository interface {
	ReadOnlyRepository
	WriteRepository
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

func (r *repository) GetTransactions(ctx context.Context, address string) ([]*blockchain.Transaction, error) {
	return r.storage.Get(ctx, address)
}

func (r *repository) InsertTransactions(ctx context.Context, addressTransactions []*AddressTransaction) error {
	data := make(map[string][]*blockchain.Transaction)
	for _, addressTransaction := range addressTransactions {
		data[addressTransaction.ID.String()] = append(data[addressTransaction.ID.String()], addressTransaction.Transaction)
	}
	return r.storage.InsertBatch(ctx, data)
}

// AddressTransactionID is a unique identifier for address transaction
// it can be more complex, but for the sake of simplicity, I will use only address
type AddressTransactionID string

func NewAddressTransactionID(
	address string,
) AddressTransactionID {
	return AddressTransactionID(
		fmt.Sprintf("%s", address),
	)
}

func (a AddressTransactionID) String() string {
	return string(a)
}

// AddressTransaction is a representation of address transaction
// It will be converted to some model representation of transaction and stored in storage
type AddressTransaction struct {
	ID          AddressTransactionID    `json:"id"`
	Transaction *blockchain.Transaction `json:"transaction"`
}
