package main

import (
	"context"
	"github.com/veljkomatic/be-homework/cmd/parser-service/internal/server"
	"github.com/veljkomatic/be-homework/pkg/parser"
	"github.com/veljkomatic/be-homework/pkg/storage/block"
	"github.com/veljkomatic/be-homework/pkg/storage/transaction"
	"log"
	"time"

	processor "github.com/veljkomatic/be-homework/cmd/parser-service/internal/block_processor"
	filter "github.com/veljkomatic/be-homework/cmd/parser-service/internal/transaction_filter"
	"github.com/veljkomatic/be-homework/pkg/blockchain"
	"github.com/veljkomatic/be-homework/pkg/provider"
	subscriberpkg "github.com/veljkomatic/be-homework/pkg/subscriber"
)

const (
	bufferSize        = 100
	heartbeatInterval = 5 * time.Minute
	serverPort        = "8080"
)

func main() {
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	app := &App{}
	app.init()
	defer app.close(ctx)

	app.startProcessing(ctx)
	app.startServer()

	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("Context finished:", ctx.Err())
			return
		case <-heartbeatTicker.C:
			log.Println("Heartbeat")
		}
	}
}

// App is the main application
type App struct {
	blockRepository       block.Repository
	transactionRepository transaction.Repository
	subscriber            subscriberpkg.Subscriber

	processedBlockChannel chan *blockchain.Block
	blockProcessor        processor.BlockProcessor
	transactionFilter     filter.TransactionFilter
}

// init initializes the application
func (a *App) init() {
	a.initRepositories()
	a.initChannels()
	a.initSubscriber()
	a.initBlockProcessor()
	a.initTransactionFilter()
}

// close closes the application
func (a *App) close(ctx context.Context) {
	a.blockProcessor.Close(ctx)
	a.transactionFilter.Close(ctx)
}

// startProcessing starts the processing of new blocks and transactions
func (a *App) startProcessing(ctx context.Context) {
	go a.blockProcessor.Start(ctx)
	go a.blockProcessor.HandleFailedBlocks(ctx)
	go a.transactionFilter.Listen(ctx)
}

// startServer starts the rest server
func (a *App) startServer() {
	parser := parser.NewParser(a.subscriber, a.transactionRepository, a.blockRepository)
	service := server.NewService(parser)
	go server.StartServer(service, serverPort)
}

// initRepositories initializes the repositories
func (a *App) initRepositories() {
	a.blockRepository = block.NewRepository(block.NewStorage())
	a.transactionRepository = transaction.NewRepository(transaction.NewStorage())
}

// initChannels initializes the channels
func (a *App) initChannels() {
	a.processedBlockChannel = make(chan *blockchain.Block, bufferSize)
}

func (a *App) initSubscriber() {
	a.subscriber = subscriberpkg.NewSubscriber()
}

// initBlockProcessor initializes the block processor
func (a *App) initBlockProcessor() {
	rpcProvider := provider.NewProvider()
	failedToProcessChan := make(chan *blockchain.BlockNumber, bufferSize)
	a.blockProcessor = processor.NewBlockProcessor(rpcProvider, a.blockRepository, a.processedBlockChannel, failedToProcessChan)
}

// initTransactionFilter initializes the transaction filter
func (a *App) initTransactionFilter() {
	subscriptionFilter := subscriberpkg.NewFilter(a.subscriber)
	a.transactionFilter = filter.NewTransactionFilter(subscriptionFilter, a.processedBlockChannel, a.transactionRepository)
}
