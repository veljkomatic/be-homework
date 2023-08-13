package parser

import (
	"context"
	"github.com/veljkomatic/be-homework/pkg/blockchain"
	"github.com/veljkomatic/be-homework/pkg/storage/block"
	"github.com/veljkomatic/be-homework/pkg/storage/transaction"
	subscriberpkg "github.com/veljkomatic/be-homework/pkg/subscriber"
	"log"
)

type Parser interface {
	// GetCurrentBlock returns last processed block number
	GetCurrentBlock(ctx context.Context) int

	// Subscribe add address to observer
	Subscribe(ctx context.Context, address string) bool

	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(ctx context.Context, address string) []*blockchain.Transaction
}

type parser struct {
	subscriber            subscriberpkg.Subscriber
	transactionRepository transaction.Repository
	blockRepository       block.Repository
}

var _ Parser = (*parser)(nil)

func NewParser(
	subscriber subscriberpkg.Subscriber,
	transactionRepository transaction.Repository,
	blockRepository block.Repository,
) Parser {
	return &parser{
		subscriber:            subscriber,
		transactionRepository: transactionRepository,
		blockRepository:       blockRepository,
	}
}

func (p *parser) GetCurrentBlock(ctx context.Context) int {
	currentBlockNumber, err := p.blockRepository.GetCurrentBlockNumber(ctx)
	if err != nil {
		log.Println("error getting current block number", err)
		return 0
	}
	return int(currentBlockNumber)
}

func (p *parser) Subscribe(ctx context.Context, address string) bool {
	err := p.subscriber.Subscribe(ctx, address)
	if err != nil {
		log.Println("error subscribing to address", address, err)
		return false
	}
	return true
}

func (p *parser) GetTransactions(ctx context.Context, address string) []*blockchain.Transaction {
	txs, err := p.transactionRepository.GetTransactions(ctx, address)
	if err != nil {
		log.Println("error getting transactions for address", address, err)
		return nil
	}
	return txs
}
