package server

import (
	"context"
	"github.com/veljkomatic/be-homework/pkg/blockchain"
	"github.com/veljkomatic/be-homework/pkg/parser"
)

type Service interface {
	GetCurrentBlockNumber(ctx context.Context) int
	Subscribe(ctx context.Context, address string) bool
	GetTransactions(ctx context.Context, address string) []*blockchain.Transaction
}

var _ Service = (*service)(nil)

type service struct {
	parser parser.Parser
}

func NewService(parser parser.Parser) Service {
	return &service{
		parser: parser,
	}
}

func (s *service) GetCurrentBlockNumber(ctx context.Context) int {
	return s.parser.GetCurrentBlock(ctx)
}

func (s *service) Subscribe(ctx context.Context, address string) bool {
	return s.parser.Subscribe(ctx, address)
}

func (s *service) GetTransactions(ctx context.Context, address string) []*blockchain.Transaction {
	return s.parser.GetTransactions(ctx, address)
}
