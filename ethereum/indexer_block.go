package ethereum

import (
	"context"
	"log"
)

// indexBlockProcess indexes the blocks with the numbers from the given channel.
func (i *Indexer) indexBlockProcess(ctx context.Context, blockCh <-chan uint64) error {
	for {
		select {
		case number, ok := <-blockCh:
			log.Printf("indexing block %d...", number)
			if !ok {
				return nil
			}
			block, err := i.rpc.GetBlockByNumber(ctx, number, true)
			if err != nil {
				log.Printf("failed to get block: %v", err)
				continue
			}
			if err := i.indexBlock(ctx, block); err != nil {
				log.Printf("failed to index block: %v", err)
				continue
			}
			log.Printf("successfully indexed block %d", number)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
