package block_processor

import (
	"context"
	"github.com/veljkomatic/be-homework/pkg/storage/block"
	"log"
	"sync"
	"time"

	"github.com/veljkomatic/be-homework/pkg/blockchain"
	"github.com/veljkomatic/be-homework/pkg/provider"
)

const (
	blockMonitorInterval         = 10 * time.Second
	retryDelay                   = 2 * time.Second
	maxRetries                   = 3
	maxConcurrentBlocksToProcess = 20
	maxConcurrentBlockRetries    = 10
)

// BlockProcessor is responsible for processing new blocks
type BlockProcessor interface {
	// Start starts the block processor
	Start(ctx context.Context)
	// HandleFailedBlocks handles process failed blocks
	HandleFailedBlocks(ctx context.Context)
	// Close closes the block processor
	Close(ctx context.Context)
}

type blockProcessor struct {
	rpcProvider     provider.Provider
	blockRepository block.Repository

	processedBlockChan  chan<- *blockchain.Block
	failedToProcessChan chan *blockchain.BlockNumber
	processingMutex     sync.Mutex
}

func NewBlockProcessor(
	rpcProvider provider.Provider,
	blockRepository block.Repository,

	processedBlockChannel chan<- *blockchain.Block,
	failedToProcessChan chan *blockchain.BlockNumber,
) BlockProcessor {
	return &blockProcessor{
		rpcProvider:     rpcProvider,
		blockRepository: blockRepository,

		processedBlockChan:  processedBlockChannel,
		failedToProcessChan: failedToProcessChan,
	}
}

func (p *blockProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(blockMonitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := p.processNewBlocks(ctx)
			if err != nil {
				log.Println(ctx, err, "process new blocks")
			}
		}
	}
}

func (p *blockProcessor) HandleFailedBlocks(ctx context.Context) {
	// The semaphore channel
	retrySemaphore := make(chan struct{}, maxConcurrentBlockRetries)

	for {
		select {
		case <-ctx.Done():
			close(retrySemaphore)
			return
		case block := <-p.failedToProcessChan:
			// this should never happen
			// defensive programming
			if block == nil {
				continue
			}
			// Handle the failed block.
			log.Printf("Block %d failed. Retrying in %v...", block, retryDelay)
			time.Sleep(retryDelay)

			retrySemaphore <- struct{}{} // acquire a semaphore slot

			go func(block *blockchain.BlockNumber) {
				defer func() { <-retrySemaphore }() // release a semaphore slot when we're done
				err := p.processBlock(ctx, *block)
				if err != nil {
					log.Println(ctx, err, "Error processing block: %s.\n", err, block)
				}
			}(block)
		}
	}
}

func (p *blockProcessor) processNewBlocks(ctx context.Context) error {
	// in production this should be handled by a distributed lock
	p.processingMutex.Lock()
	defer p.processingMutex.Unlock()

	latestBlockNumber, err := p.rpcProvider.GetLatestBlockNumber(ctx)
	if err != nil {
		return err
	}
	currentBlockNumber, err := p.blockRepository.GetCurrentBlockNumber(ctx)
	if err != nil {
		return err
	}
	if currentBlockNumber == blockchain.EarliestBlockNumber {
		err := p.processBlock(ctx, latestBlockNumber)
		if err != nil {
			return err
		}
		return p.blockRepository.SaveBlockNumber(ctx, latestBlockNumber)
	}

	if latestBlockNumber > currentBlockNumber {
		// call processBlocksInParallel in separate goroutine
		// to avoid blocking the main thread
		go p.processBlocksInParallel(ctx, currentBlockNumber.Inc(), latestBlockNumber)
		// optimistic update
		// if error occurs, the block will be retried and eventually processed as success
		// if the block is not processed even in HandleFailedBlocks,
		// we could send failed blocks to queue and retry later via a separate job
		return p.blockRepository.SaveBlockNumber(ctx, latestBlockNumber)
	}

	//TODO: handle reorgs in future
	return nil
}

func (p *blockProcessor) processBlock(ctx context.Context, blockNumber blockchain.BlockNumber) error {
	var currentRetry int

	log.Printf("Processing block number %d.", blockNumber)

	var err error
	var block *blockchain.Block

	for currentRetry < maxRetries {
		block, err = p.rpcProvider.GetBlockByNumber(ctx, blockNumber)
		if err != nil {
			log.Println(ctx, err, "Error fetching block number %d: %s. Retry %d/%d.\n", blockNumber, err, currentRetry+1, maxRetries)
			currentRetry++
			continue
		}
		p.processedBlockChan <- block
		break
	}

	if currentRetry == maxRetries {
		log.Println(ctx, "Failed to fetch block number %d after %d retries.\n", blockNumber, maxRetries)
		return err
	}
	return nil
}

func (p *blockProcessor) processBlocksInParallel(ctx context.Context, startBlockNumber blockchain.BlockNumber, endBlockNumber blockchain.BlockNumber) {
	// The semaphore channel
	processSemaphore := make(chan struct{}, maxConcurrentBlocksToProcess)

	wg := sync.WaitGroup{}
	for i := startBlockNumber.ToInt64(); i <= endBlockNumber.ToInt64(); i++ {
		wg.Add(1)
		processSemaphore <- struct{}{} // Acquire a semaphore slot

		go func(blockNumber blockchain.BlockNumber) {
			defer wg.Done()
			defer func() { <-processSemaphore }() // Release a semaphore slot when we're done

			err := p.processBlock(ctx, blockNumber)
			if err != nil {
				p.failedToProcessChan <- &blockNumber
			}
		}(blockchain.BlockNumber(i))
	}

	wg.Wait()
	close(processSemaphore)
}

func (p *blockProcessor) Close(ctx context.Context) {
	close(p.processedBlockChan)
	close(p.failedToProcessChan)
}
