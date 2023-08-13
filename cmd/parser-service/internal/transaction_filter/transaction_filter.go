package transaction_filter

import (
	"context"
	"github.com/veljkomatic/be-homework/pkg/blockchain"
	"github.com/veljkomatic/be-homework/pkg/storage/transaction"
	"github.com/veljkomatic/be-homework/pkg/subscriber"
	"log"
)

const maxConcurrentFilters = 20

// TransactionFilter is a service that listens for processed new blocks and filters transactions.
type TransactionFilter interface {
	// Listen starts listening for new blocks and filters transactions.
	Listen(ctx context.Context)
	// Close closes the transaction filter.
	Close(ctx context.Context)
}

type transactionFilter struct {
	filter                subscriber.Filter
	processedBlockChannel <-chan *blockchain.Block
	transactionRepository transaction.WriteRepository
}

func NewTransactionFilter(
	filter subscriber.Filter,
	processedBlockChannel <-chan *blockchain.Block,
	transactionRepository transaction.WriteRepository,
) TransactionFilter {
	return &transactionFilter{
		filter:                filter,
		processedBlockChannel: processedBlockChannel,
		transactionRepository: transactionRepository,
	}
}

func (t *transactionFilter) Listen(ctx context.Context) {
	// semaphore to limit the number of concurrent filters
	filterSemaphore := make(chan struct{}, maxConcurrentFilters)

	for {
		select {
		case <-ctx.Done():
			close(filterSemaphore)
			return
		case block := <-t.processedBlockChannel:
			// defensive programming
			// we should never receive a nil block
			if block == nil {
				continue
			}
			filterSemaphore <- struct{}{} // acquire a semaphore slot
			// filter transactions in a separate goroutine
			// so we can continue listening for new blocks
			go func(block *blockchain.Block) {
				defer func() { <-filterSemaphore }() // release the semaphore slot
				t.filterTransactions(ctx, block)
			}(block)
		}
	}
}

// filterTransactions filters transactions from a block if they match the filter and stores them in the database.
// here we are using a simple filter that checks if the transaction's from or to address matches the filter.
// in a real world scenario we would probably want to use a bloom filter to check if the transaction's from or to address matches the filter.
// here we could send filtered transactions to a queue so notification service can send notifications to subscribers.
func (t *transactionFilter) filterTransactions(ctx context.Context, block *blockchain.Block) {
	filteredTransactions := make([]*transaction.AddressTransaction, 0, len(block.Transactions))
	for _, tx := range block.Transactions {
		if t.filter.Test(ctx, tx.From) {
			filteredTransactions = append(filteredTransactions, &transaction.AddressTransaction{
				ID:          transaction.NewAddressTransactionID(tx.From),
				Transaction: tx,
			})
		}
		if t.filter.Test(ctx, tx.To) {
			filteredTransactions = append(filteredTransactions, &transaction.AddressTransaction{
				ID:          transaction.NewAddressTransactionID(tx.To),
				Transaction: tx,
			})
		}
	}
	if err := t.storeObservedTransactions(ctx, filteredTransactions); err != nil {
		log.Println(ctx, err, "Error storing observed transactions")
	}
}

// storeObservedTransactions stores filtered transactions in the database (in memory).
// in a real world database would be a persistent storage, some NoSQL database like MongoDB or Cassandra.
func (t *transactionFilter) storeObservedTransactions(ctx context.Context, filteredTransactions []*transaction.AddressTransaction) error {
	if len(filteredTransactions) == 0 {
		return nil
	}

	var currentRetry int
	const maxRetries = 3

	var err error

	for currentRetry < maxRetries {
		if err := t.transactionRepository.InsertTransactions(ctx, filteredTransactions); err != nil {
			log.Println(ctx, err, "Error inserting transactions: %s. Retry %d/%d.\n", err, currentRetry+1, maxRetries)
			currentRetry++
			continue
		}
		break
	}

	if currentRetry == maxRetries {
		log.Println(ctx, err, "Error inserting transactions: %s. Max retries exceeded.\n", err)
		return err
	}
	return nil
}

func (t *transactionFilter) Close(ctx context.Context) {}
