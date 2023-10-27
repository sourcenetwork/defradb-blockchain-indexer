package evm

import (
	"context"
	"log"
	"time"
)

// blockTime is the time between new blocks
var blockTime = 12 * time.Second

// indexHeadProcess adds block numbers to the given channel every time a new block is added.
func (i *Indexer) indexHeadProcess(ctx context.Context, blockCh chan<- uint64) error {
	var head uint64
	for {
		select {
		case <-time.Tick(blockTime):
			number, err := i.rpc.BlockNumber(ctx)
			if err != nil {
				log.Printf("failed to get block number: %v", err)
				continue
			}
			if head == 0 && number > 0 {
				head = number - 1
			}
			for i := head; i < number; i++ {
				blockCh <- number
			}
			head = number
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
