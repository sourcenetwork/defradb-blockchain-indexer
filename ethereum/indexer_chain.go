package ethereum

import (
	"context"
	"log"
	"time"
)

// blockTime is the time between new blocks
var blockTime = 12 * time.Second

// indexChainProcess adds block numbers to the given channel every time a new block is added.
func (i *Indexer) indexChainProcess(ctx context.Context, blockCh chan<- uint64, head *uint64) error {
	for {
		select {
		case <-time.Tick(blockTime):
			number, err := i.rpc.BlockNumber(ctx)
			if err != nil {
				log.Printf("failed to get block number: %v", err)
				continue
			}
			if head == nil {
				head = &number
			}
			for i := *head; i <= number; i++ {
				blockCh <- number
			}
			next := number + 1
			head = &next
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
