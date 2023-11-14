package ethereum

import (
	"context"
	"log"
)

// indexHeaderProcess indexes the headers with the numbers from the given channel.
func (i *Indexer) indexHeaderProcess(ctx context.Context, blockCh <-chan uint64) error {
	for {
		select {
		case number, ok := <-blockCh:
			log.Printf("indexing header %d...", number)
			if !ok {
				return nil
			}
			header, err := i.rpc.GetBlockByNumber(ctx, number, false)
			if err != nil {
				log.Printf("failed to get header: %v", err)
				continue
			}
			if err := i.indexHeader(ctx, header); err != nil {
				log.Printf("failed to index header: %v", err)
				continue
			}
			log.Printf("successfully indexed header %d", number)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
