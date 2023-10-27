package evm

import (
	"context"
	_ "embed"
	"log"

	"github.com/sourcenetwork/defradb/client"
)

//go:embed schema.gql
var schema string

type Indexer struct {
	db  client.DB
	rpc *RpcClient
}

func NewIndexer(ctx context.Context, db client.DB, url string) (*Indexer, error) {
	_, err := db.AddSchema(ctx, schema)
	if err != nil {
		return nil, err
	}
	return &Indexer{
		db:  db,
		rpc: NewRpcClient(url),
	}, nil
}

// Start runs the indexer with the given number of processes.
func (idx *Indexer) Start(ctx context.Context, count int) error {
	blockCh := make(chan uint64)
	defer close(blockCh)

	log.Printf("Starting indexer with %d processes...", count)
	for i := 0; i < count; i++ {
		go func() { idx.indexBlockProcess(ctx, blockCh) }()
	}
	return idx.indexHeadProcess(ctx, blockCh)
}

func (idx *Indexer) indexBlock(ctx context.Context, block *Block) error {
	txn, err := idx.db.NewTxn(ctx, false)
	defer txn.Discard(ctx)

	blocks, err := idx.db.GetCollectionByName(ctx, "Block")
	if err != nil {
		return err
	}
	transactions, err := idx.db.GetCollectionByName(ctx, "Transaction")
	if err != nil {
		return err
	}

	blockDoc, err := client.NewDocFromMap(block.ToMap())
	if err != nil {
		return err
	}
	err = blocks.WithTxn(txn).Create(ctx, blockDoc)
	if err != nil {
		return err
	}
	for _, transaction := range block.Transactions {
		transactionDoc, err := client.NewDocFromMap(transaction.ToMap())
		if err != nil {
			return err
		}
		err = transactions.WithTxn(txn).Create(ctx, transactionDoc)
		if err != nil {
			return err
		}
	}
	return txn.Commit(ctx)
}
